import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.18",
    settings: {
      // istanbul upgrade occurred before the london jinxfork, so is compatible with fury's evm
      evmVersion: "istanbul",
      // optimize build for deployment to mainnet!
      optimizer: {
        enabled: true,
        runs: 1000,
      },
    },
  },
  networks: {
    // futool's local network
    highbury: {
      url: "http://34.93.171.55:8545",
      accounts: [
        "8BC170FE21C3137F4058973E47CF3FC45070CAB6C3D4DAB63FBE8CAA5A846BBB",
      ],
    },
    protonet: {
      url: "https://evm.app.protonet.us-east.production.fury.io:443",
      accounts: [
        "247069F0BC3A5914CB2FD41E4133BBDAA6DBED9F47A01B9F110B5602C6E4CDD9",
      ],
    },
    internal_testnet: {
      url: "https://evm.data.internal.testnet.us-east.production.fury.io:443",
      accounts: [
        "247069F0BC3A5914CB2FD41E4133BBDAA6DBED9F47A01B9F110B5602C6E4CDD9",
      ],
    },
  },
};

export default config;
