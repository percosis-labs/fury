/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { Bridge, BridgeInterface } from "../Bridge";

const _abi = [
  {
    inputs: [
      {
        internalType: "address",
        name: "relayer_",
        type: "address",
      },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "token",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "toFuryAddr",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "lockSequence",
        type: "uint256",
      },
    ],
    name: "Lock",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "token",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "toAddr",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "unlockSequence",
        type: "uint256",
      },
    ],
    name: "Unlock",
    type: "event",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "token",
        type: "address",
      },
      {
        internalType: "address",
        name: "toFuryAddr",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
    ],
    name: "lock",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "relayer",
    outputs: [
      {
        internalType: "address",
        name: "",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "token",
        type: "address",
      },
      {
        internalType: "address",
        name: "toAddr",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "amount",
        type: "uint256",
      },
      {
        internalType: "uint256",
        name: "unlockSequence",
        type: "uint256",
      },
    ],
    name: "unlock",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
];

const _bytecode =
  "0x608060405234801561001057600080fd5b50604051610e3e380380610e3e833981810160405281019061003291906100ee565b60006001600081905550806001819055505080600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061011b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100bb82610090565b9050919050565b6100cb816100b0565b81146100d657600080fd5b50565b6000815190506100e8816100c2565b92915050565b6000602082840312156101045761010361008b565b5b6000610112848285016100d9565b91505092915050565b610d148061012a6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634202e907146100465780637750c9f0146100625780638406c0791461007e575b600080fd5b610060600480360381019061005b919061077c565b61009c565b005b61007c600480360381019061007791906107e3565b6101d4565b005b6100866102a3565b6040516100939190610845565b60405180910390f35b6100a46102cd565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610134576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161012b906108bd565b60405180910390fd5b61015f83838673ffffffffffffffffffffffffffffffffffffffff1661031d9092919063ffffffff16565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fc1640cf787ea538af4a68163da63c6da8b4577194278080f5b040a75df00389984846040516101be9291906108ec565b60405180910390a36101ce6103a3565b50505050565b6101dc6102cd565b6101e46103ad565b6102113330838673ffffffffffffffffffffffffffffffffffffffff166103ba909392919063ffffffff16565b8173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f749e347e95185169edffb86e003abcbf08f3510641663b00600acf707f098f4784610280610443565b60405161028e9291906108ec565b60405180910390a461029e6103a3565b505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60026000541415610313576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161030a90610961565b60405180910390fd5b6002600081905550565b61039e8363a9059cbb60e01b848460405160240161033c929190610981565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061044d565b505050565b6001600081905550565b6001805401600181905550565b61043d846323b872dd60e01b8585856040516024016103db939291906109aa565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061044d565b50505050565b6000600154905090565b60006104af826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166105159092919063ffffffff16565b90506000815114806104d15750808060200190518101906104d09190610a19565b5b610510576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050790610ab8565b60405180910390fd5b505050565b6060610524848460008561052d565b90509392505050565b606082471015610572576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161056990610b4a565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161059b9190610be4565b60006040518083038185875af1925050503d80600081146105d8576040519150601f19603f3d011682016040523d82523d6000602084013e6105dd565b606091505b50915091506105ee878383876105fa565b92505050949350505050565b6060831561065d576000835114156106555761061585610670565b610654576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064b90610c47565b60405180910390fd5b5b829050610668565b6106678383610693565b5b949350505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6000825111156106a65781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106da9190610cbc565b60405180910390fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610713826106e8565b9050919050565b61072381610708565b811461072e57600080fd5b50565b6000813590506107408161071a565b92915050565b6000819050919050565b61075981610746565b811461076457600080fd5b50565b60008135905061077681610750565b92915050565b60008060008060808587031215610796576107956106e3565b5b60006107a487828801610731565b94505060206107b587828801610731565b93505060406107c687828801610767565b92505060606107d787828801610767565b91505092959194509250565b6000806000606084860312156107fc576107fb6106e3565b5b600061080a86828701610731565b935050602061081b86828701610731565b925050604061082c86828701610767565b9150509250925092565b61083f81610708565b82525050565b600060208201905061085a6000830184610836565b92915050565b600082825260208201905092915050565b7f4272696467653a20756e74727573746564206164647265737300000000000000600082015250565b60006108a7601983610860565b91506108b282610871565b602082019050919050565b600060208201905081810360008301526108d68161089a565b9050919050565b6108e681610746565b82525050565b600060408201905061090160008301856108dd565b61090e60208301846108dd565b9392505050565b7f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00600082015250565b600061094b601f83610860565b915061095682610915565b602082019050919050565b6000602082019050818103600083015261097a8161093e565b9050919050565b60006040820190506109966000830185610836565b6109a360208301846108dd565b9392505050565b60006060820190506109bf6000830186610836565b6109cc6020830185610836565b6109d960408301846108dd565b949350505050565b60008115159050919050565b6109f6816109e1565b8114610a0157600080fd5b50565b600081519050610a13816109ed565b92915050565b600060208284031215610a2f57610a2e6106e3565b5b6000610a3d84828501610a04565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000610aa2602a83610860565b9150610aad82610a46565b604082019050919050565b60006020820190508181036000830152610ad181610a95565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b6000610b34602683610860565b9150610b3f82610ad8565b604082019050919050565b60006020820190508181036000830152610b6381610b27565b9050919050565b600081519050919050565b600081905092915050565b60005b83811015610b9e578082015181840152602081019050610b83565b83811115610bad576000848401525b50505050565b6000610bbe82610b6a565b610bc88185610b75565b9350610bd8818560208601610b80565b80840191505092915050565b6000610bf08284610bb3565b915081905092915050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b6000610c31601d83610860565b9150610c3c82610bfb565b602082019050919050565b60006020820190508181036000830152610c6081610c24565b9050919050565b600081519050919050565b6000601f19601f8301169050919050565b6000610c8e82610c67565b610c988185610860565b9350610ca8818560208601610b80565b610cb181610c72565b840191505092915050565b60006020820190508181036000830152610cd68184610c83565b90509291505056fea2646970667358221220d9797a59ee85106033b750f20c7647c040d4c55daa05b1f92ed65498dfbc52fa64736f6c63430008090033";

type BridgeConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: BridgeConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class Bridge__factory extends ContractFactory {
  constructor(...args: BridgeConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
    this.contractName = "Bridge";
  }

  deploy(
    relayer_: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<Bridge> {
    return super.deploy(relayer_, overrides || {}) as Promise<Bridge>;
  }
  getDeployTransaction(
    relayer_: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(relayer_, overrides || {});
  }
  attach(address: string): Bridge {
    return super.attach(address) as Bridge;
  }
  connect(signer: Signer): Bridge__factory {
    return super.connect(signer) as Bridge__factory;
  }
  static readonly contractName: "Bridge";
  public readonly contractName: "Bridge";
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): BridgeInterface {
    return new utils.Interface(_abi) as BridgeInterface;
  }
  static connect(address: string, signerOrProvider: Signer | Provider): Bridge {
    return new Contract(address, _abi, signerOrProvider) as Bridge;
  }
}
