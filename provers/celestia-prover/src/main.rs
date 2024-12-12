use tonic::{transport::Server, Request, Response, Status};
use std::fs;
use ibc_client_tendermint_types::ConsensusState;
// Import the generated proto rust code
pub mod prover {
    tonic::include_proto!("celestia.prover.v1");
}

use prover::prover_server::{Prover, ProverServer};
use prover::{
    ProveStateTransitionRequest, ProveStateTransitionResponse,
    ProveMembershipRequest, ProveMembershipResponse,
};
use sp1_ics07_tendermint_prover::{
    programs::UpdateClientProgram,
    prover::{SP1ICS07TendermintProver, SupportedProofType},
};
use tendermint_rpc::HttpClient;

#[derive(Default)]
pub struct ProverService {
    tendermint_prover: SP1ICS07TendermintProver<UpdateClientProgram>,
    tendermint_rpc_client: HttpClient,
    membership_prover: SP1MembershipProver<MembershipProofProgram>,
    contract: sp1_ics07_tendermint::Contract,
}

#[tonic::async_trait]
impl Prover for ProverService {
    fn new() -> Self {
        let rpc_url = env::var("EVM_RPC_URL").expect("EVM_RPC_URL not set");
        let contract_address = env::var("EVM_CONTRACT_ADDRESS").expect("EVM_CONTRACT_ADDRESS not set");

        let provider = ProviderBuilder::new()
            .with_recommended_fillers()
            .on_http(Url::parse(rpc_url.as_str())?);
        let contract = sp1_ics07_tendermint::new(contract_address.parse()?, provider);

        Self {
            tendermint_prover: SP1ICS07TendermintProver::new(SupportedProofType::Groth16),
            tendermint_rpc_client: HttpClient::from_env(),
            membership_prover: SP1MembershipProver::new(SupportedProofType::Groth16),
            contract: contract,
        }
    }

    async fn prove_state_transition(
        &self,
        request: Request<ProveStateTransitionRequest>,
    ) -> Result<Response<ProveStateTransitionResponse>, Status> {
        println!("Got state transition request: {:?}", request);
        let client_state = self.contract.getClientState().call().await?._0;
        // fetch the light block at the latest height of the client state
        let trusted_light_block = self.tendermint_rpc_client.get_light_block(Some(client_state.latestHeight.revisionHeight)).await?;
        // fetch the latest light block
        let target_light_block = self.tendermint_rpc_client.get_light_block(None).await?;

        let trusted_consensus_state = trusted_light_block.to_consensus_state().into();
        let proposed_header = target_light_block.into_header(&trusted_light_block);

        let now = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)?
            .as_secs();

        let proof = self.prover.generate_proof(
            &client_state,
            &trusted_consensus_state,
            &proposed_header,
            now,
        );
        
        // Implement your state transition proof logic here
        let response = ProveStateTransitionResponse {
            proof: proof,
        };

        Ok(Response::new(response))
    }

    async fn prove_membership(
        &self,
        request: Request<ProveMembershipRequest>,
    ) -> Result<Response<ProveMembershipResponse>, Status> {
        println!("Got membership request: {:?}", request);

        // Implement your membership proof logic here
        let response = ProveMembershipResponse {
            proof: vec![], // Replace with actual proof
            height: 0,     // Replace with actual height
        };

        Ok(Response::new(response))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let prover = ProverService::default();

    println!("Prover Server listening on {}", addr);

    // Load the file descriptor set
    let file_descriptor_set = fs::read("proto_descriptor.bin")?;

    Server::builder()
        .add_service(ProverServer::new(prover))
        .add_service(
            tonic_reflection::server::Builder::configure()
                .register_encoded_file_descriptor_set(&file_descriptor_set)
                .build()
                .unwrap()
        )
        .serve(addr)
        .await?;

    Ok(())
}