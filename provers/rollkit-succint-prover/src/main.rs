use std::sync::Arc;
use tendermint_rpc::HttpClient;
use tokio::sync::mpsc;
use tokio::time::{sleep, Duration};
use tonic;

// Prover service client implementation
pub mod prover {
    tonic::include_proto!("celestia.prover.v1");
}

use prover::prover_client::ProverClient;
use prover::{KeyValuePair, ProveMembershipRequest, ProveStateTransitionRequest};

// Configuration struct for the relayer
#[derive(Clone, Debug)]
pub struct RelayerConfig {
    celestia_rpc: String,
    rollup_rpc: String,
    prover_endpoint: String,
    celestia_ibc_contract: String,
}

// Event types from the rollup
#[derive(Debug)]
enum RollupEvent {
    IBCPacket {
        sequence: u64,
        source_port: String,
        source_channel: String,
        data: Vec<u8>,
        timeout_height: u64,
    },
    StateUpdate {
        height: u64,
        state_root: Vec<u8>,
    },
}

// Main relayer struct
pub struct Relayer {
    config: RelayerConfig,
    celestia_client: HttpClient,
    prover_client: ProverClient<tonic::transport::Channel>,
}

impl Relayer {
    pub async fn new(config: RelayerConfig) -> Result<Self, Box<dyn std::error::Error>> {
        let celestia_client = HttpClient::new(config.celestia_rpc.as_str())?;
        let prover_client = ProverClient::connect(config.prover_endpoint.clone()).await?;

        Ok(Self {
            config,
            celestia_client,
            prover_client,
        })
    }

    // Start listening for events from the rollup
    pub async fn start(&self) -> Result<(), Box<dyn std::error::Error>> {
        let (tx, mut rx) = mpsc::channel(100);
        let tx = Arc::new(tx);

        // Spawn rollup event listener
        let rollup_tx = tx.clone();
        let rollup_endpoint = self.config.rollup_rpc.clone();
        tokio::spawn(async move {
            Self::listen_rollup_events(rollup_endpoint, rollup_tx).await;
        });

        // Process events
        while let Some(event) = rx.recv().await {
            self.handle_event(event).await?;
        }

        Ok(())
    }

    async fn listen_rollup_events(
        endpoint: String,
        tx: Arc<mpsc::Sender<RollupEvent>>,
    ) -> Result<(), Box<dyn std::error::Error>> {
        // TODO: Setup websocket connection to rollup
        // Subscribe to:
        // 1. IBC packet events
        // 2. State update events
        loop {
            sleep(Duration::from_secs(1)).await;
        }
    }

    async fn handle_event(&self, event: RollupEvent) -> Result<(), Box<dyn std::error::Error>> {
        match event {
            RollupEvent::IBCPacket {
                sequence,
                source_port,
                source_channel,
                data,
                timeout_height,
            } => {
                self.handle_packet(sequence, source_port, source_channel, data, timeout_height)
                    .await?;
            }
            RollupEvent::StateUpdate { height, state_root } => {
                self.handle_state_update(height, state_root).await?;
            }
        }
        Ok(())
    }

    async fn handle_packet(
        &self,
        sequence: u64,
        source_port: String,
        source_channel: String,
        data: Vec<u8>,
        timeout_height: u64,
    ) -> Result<(), Box<dyn std::error::Error>> {
        // 1. Get membership proof for the packet
        let membership_proof = self.get_membership_proof(sequence, &data).await?;

        // 2. Submit to Celestia
        self.submit_to_celestia(sequence, &data, &membership_proof)
            .await?;

        Ok(())
    }

    async fn handle_state_update(
        &self,
        height: u64,
        state_root: Vec<u8>,
    ) -> Result<(), Box<dyn std::error::Error>> {
        // Get state transition proof
        let proof = self.get_state_transition_proof(height, &state_root).await?;

        // Submit state update to Celestia
        self.submit_state_update(height, &state_root, &proof)
            .await?;

        Ok(())
    }

    async fn get_membership_proof(
        &self,
        sequence: u64,
        data: &[u8],
    ) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
        let request = tonic::Request::new(ProveMembershipRequest {
            height: sequence as i64,
            key_value_pairs: [].to_vec(),
        });

        let response = self.prover_client.clone().prove_membership(request).await?;

        Ok(response.into_inner().proof)
    }

    async fn get_state_transition_proof(
        &self,
        height: u64,
        state_root: &[u8],
    ) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
        let request = tonic::Request::new(ProveStateTransitionRequest {});

        let response = self
            .prover_client
            .clone()
            .prove_state_transition(request)
            .await?;

        Ok(response.into_inner().proof)
    }

    async fn submit_to_celestia(
        &self,
        sequence: u64,
        data: &[u8],
        proof: &[u8],
    ) -> Result<(), Box<dyn std::error::Error>> {
        // TODO: Submit IBC packet with proof to Celestia
        // 1. Format transaction
        // 2. Submit to Celestia RPC
        // 3. Wait for confirmation
        Ok(())
    }

    async fn submit_state_update(
        &self,
        height: u64,
        state_root: &[u8],
        proof: &[u8],
    ) -> Result<(), Box<dyn std::error::Error>> {
        // TODO: Submit state update with proof to Celestia
        // 1. Format transaction
        // 2. Submit to Celestia RPC
        // 3. Wait for confirmation
        Ok(())
    }
}

// Custom error type
#[derive(Debug)]
pub enum RelayerError {
    RPCError(String),
    ProofError(String),
    PacketError(String),
}

impl std::fmt::Display for RelayerError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            RelayerError::RPCError(msg) => write!(f, "RPC error: {}", msg),
            RelayerError::ProofError(msg) => write!(f, "Proof generation failed: {}", msg),
            RelayerError::PacketError(msg) => write!(f, "Packet relay failed: {}", msg),
        }
    }
}

impl std::error::Error for RelayerError {}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_relayer_initialization() {
        // TODO: Add initialization tests
    }

    #[tokio::test]
    async fn test_packet_relay() {
        // TODO: Add packet relaying tests
    }

    #[tokio::test]
    async fn test_state_update() {
        // TODO: Add state update tests
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Load configuration
    let config = RelayerConfig {
        celestia_rpc: std::env::var("CELESTIA_RPC")
            .unwrap_or_else(|_| "http://localhost:26657".into()),
        rollup_rpc: std::env::var("ROLLUP_RPC").unwrap_or_else(|_| "ws://localhost:8545".into()),
        prover_endpoint: std::env::var("PROVER_ENDPOINT")
            .unwrap_or_else(|_| "http://localhost:50051".into()),
        celestia_ibc_contract: std::env::var("CELESTIA_IBC_CONTRACT").unwrap_or_default(),
    };

    // Create and start the relayer
    let runtime = tokio::runtime::Runtime::new()?;
    runtime.block_on(async {
        let relayer = Relayer::new(config).await?;
        relayer.start().await
    })?;

    Ok(())
}
