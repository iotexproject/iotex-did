<p align="center">
  <img src="https://github.com/iotexproject/iotex-did/blob/master/did.png" width="480px">
</p>

# IoTeX DID Method Specification

## Status
This document is a work in progress draft. **Last updated at 12 Sep 2019**.

Self-sovereign identity is a user- and device-centric identity concept where an individual or organization or device is able to control its own identity attributes. Decentralized Identifiers (DIDs) are a new type of identifier for verifiable 'self-sovereign' digital identity for individuals, organizations and things. DIDs have the following important properties:
- Decentralized: DIDs are designed to function without a central registration authority. DIDs are registered in blockchain or other decentralized network.
- Cryptographically Verifiable: DIDs are designed to be associated with cryptographic keys and the entities controlling the DID can use those keys to prove ownership.
- Non-Reassignable: DIDs should be permanent, persistent, and non-reassignable.
- Resolvable: DIDs are made useful through resolution.

This document provides specification of how DID works within IoTeX Network. In IoTeX DID system, we allow each manufacture or namespace to store and manage DIDs through its self-managed DID smart contract. Every self-managed smart contract has to implement the SelfManagedDID interface defined as follows:

```
interface SelfManagedDID {
    function createDID(string id, bytes32 hash, string uri) public returns (string);
    function deleteDID(string did) public;
    function updateHash(string did, bytes32 hash) public;
    function updateURI(string did, string uri) public;
    function getHash(string did) public view returns (bytes32);
    function getURI(string did) public view returns (string);
}
```

The following part of this document specifies the IoTeX implementation of the SelfManagedDID interface with [DID Method](https://w3c-ccg.github.io/did-spec/#specific-did-method-schemes) [`did:io`].

## Method Name

We use the `iotex` to be our method name and a formal DID using this method need begin with following prefix: `did:io` . Furthermore, all the characters in the prefix need to be in lowercase. The string after prefix is the unique IoTeX account address of the registered user/device.

## Generate a unique ID string

Every user or device needs to register its own DID by uploading its UUID under a namespace, DID document hash, and DID document URI. Then a unique `idstring` is created as follows:

1.  Construct the `iotex` method prefix, i.e. `did:io`.
2.  Append the UUID (IoTeX account address) to the prefix. The UUID must match caller's address. 
3. The provided DID document hash and URI would be stored in the contract along the unique `idstring`.

The DID structure looks like follows:
```
struct DID {
        bool exist;
        bytes32 hash;
        string uri;
}
```

### Example

An example IoTeX DID:

```
did:io:0x5576E95935366Ebd2637D9171E4C92e60598be10
```

## CRUD Operations
### Create/Register
IoTeX DID contract implements the functionality of DID registrations. It associates every IoTeX account with a unique dencentralized identifier and a corresponding DID document. Every IoTeX account can only register a DID for itself. It is required that the registerer provides the DID document URI where its document can be accessed and the current document hash that represents the initial state of the document.
Here is a draft DID registry implementation:

```
// For reference only and subject to change

string constant didPrefix = "did:io:";
struct DID {
    bool exist;
    bytes32 hash;
    string uri;
}
mapping(string => DID) dids;

function createDID(string id, bytes32 hash, string uri) public returns (string) {
    if (bytes(id).length > 0) {
        require(compareStrings(id, addrToString(msg.sender)), "id does not match creator");
    }
    string memory resultDID = generateDIDString();
    require(!dids[resultDID].exist, "did already exists");
    dids[resultDID] = DID(true, hash, uri);
    return resultDID;
}

function generateDIDString() private view returns (string) {
    return string(abi.encodePacked(didPrefix, addrToString(msg.sender)));
}
```

An example of registration: 
```
createDID("0x5576E95935366Ebd2637D9171E4C92e60598be10", "8806157fdcbcea265623576fa72d88568db3f9ca8b36bddfe3755ae80457eaf5", "user:password@tcp(example_connection_string:3306)/")
```
Note that id argument is optional because the contract will use the caller's address anyway. Once a DID is registered, the provided DID document hash and URI would be stored in the contract along the unique DID string.

### Read
IoTeX clients can query a DID's current hash and URI given a DID string whether they own the DID or not. 
Here is a draft DID hash/URI read implementation:
```
// For reference only and subject to change

function getHash(string did) public view returns (bytes32) {
    require(dids[did].exist, "did does not exist");
    return dids[did].hash;
}

function getURI(string did) public view returns (string) {
    require(dids[did].exist, "did does not exist");
    return dids[did].uri;
}
```

### Update
Whenever a DID document is modified, the document hash will change and need to be updated in the DID contract. If a user moves the document to a new storage location, the document URI also needs to be updated in the DID contract. Every device can only update its own document associated with the DID created by itself.
Here is a draft DID hash/URI update implementation:
```
// For reference only and subject to change

modifier onlyDIDOwner(string didInput) {
    string memory didString = generateDIDString();
    if (bytes(didInput).length > 0) {
        require(compareStrings(didInput, didString), "caller does not own the given did");
    }
    require(dids[didString].exist, "caller is not a did owner");
    _;
}

function updateHash(string did, bytes32 hash) public onlyDIDOwner(did) {
    dids[generateDIDString()].hash = hash;
}

function updateURI(string did, string uri) public onlyDIDOwner(did) {
    dids[generateDIDString()].uri = uri;
}
```
Similar to the case of DID registration, did argument is optional for the update functions as well because the contract would use caller's address to derive the DID string.

## Delete
A registered DID can be deactivated anytime as long as the caller's address matches the address within the DID string. Once a DID is deactivated, the metadata of the corresponding document cannot be updated.
Here is a draft DID hash/URI delete implementation:
```
// For reference only and subject to change

function deleteDID(string did) public onlyDIDOwner(did) {
    dids[generateDIDString()].exist = false;
}
```
Similar to previous cases, the contract checks the permission to ensure that only the DID creator can deactivate its own DID. 

## DID Documentation Resolution
IoTeX DID document is a JSON-LD document containing six, optional components:
* The DID that points to the DID Document, identified by the key **id**
* A list of public keys identified by the key **publickey**
* List of protocols for authentication control of the DID and delegated capabilities identified by the key **authentication**
* A set of service endpoints that allow discovery of way to interact with the entity, identified by the key **service**
* A timestamp indicates when the DID Document was created and updated, identified by the key **created/updated**
* A digital signature for verifying the integrity of DID Document, identified by the key **proof**

Here is a draft IoTeX DID document example:

```
{
    "@context": "https://w3id.org/did/v1", 
    "id": "did:io:0x88C36867cffB66197812a9385A038cc6Dd75244b", 
    "publicKey": [{
        "id": "did:io:0x5576E95935366Ebd2637D9171E4C92e60598be10#keys-1", 
        "type": "RsaVerificationKey2018", 
        "controller": "did:io:0x56d0B5eD3D525332F00C9BC938f93598ab16AAA7", 
        "publicKeyPem": "-----BEGIN PUBLIC KEY...END PUBLIC KEY-----\r\n" }], 
    "authentication": [{ 
        "type": "RsaSignatureAuthentication2018", 
        "publicKey": "did:io:0x5576E95935366Ebd2637D9171E4C92e60598be10#keys-1" }], 
    "service": [{ 
        "id": "did:io:0x88C36867cffB66197812a9385A038cc6Dd75244b;exam_svc", 
        "type": "ExampleService", 
        "serviceEndpoint": "https://example.com/endpoint/8377464" }], 
    "created": "2018-02-08T16:03:00Z", 
    "proof": { 
        "type": "LinkedDataSignature2015", 
        "created": "2018-02-08T16:02:20Z", 
        "creator": "did:io:0x5576E95935366Ebd2637D9171E4C92e60598be10#keys-1", 
        "signatureValue": "QNB13Y7Q9...1tzjn4w==" 
    }
}    
```

## Verifiable Credentials
Verifiable claims are statements made by an entity about a 'subject' whose authorship can be cryptographically verified. They can be combined with DID documents to provide trusted service. 
IoTeX verifiable credentials include the following features:
* Claims are used to create **verifiable credentials** by issuers.
* Verifiable credentials are **decentralized** and **contextual**.
* Credential issuers decide on which claims are contained in the credentials.
* Verifiers **make their own trust decisions** about which credentials to accept.
* Verifiers do not need to contact issuers to perform verification.
* Credential holders are free to choose which credentials to carry and what information to disclose.

Here is a draft IoTeX verifiable credentials example:

```
{
    "@context": "https://w3id.org/credentials/v1", 
    "id": "did:io:0xb36A1D1778f9D5E5816682c2cC0d16C65828a6b4/credentials/1", 
    "type": ["Credential", "NameCredential"], 
    "issuer": "did:io:0x669d00D4191fB397c780212EB81B439F5Ec9967d", 
    "issued": "2019-09-01", 
    "claim": { 
        "id": "did:io:0x669d00D4191fB397c780212EB81B439F5Ec9967d", 
        "name": "John Doe", 
        "address": "..."
    },
    "proof": { 
        "type": "RsaSignature2018", 
        "created": "2017-06-18T21:19:10Z", 
        "creator": "did:io:0xb36A1D1778f9D5E5816682c2cC0d16C65828a6b4#key-1", 
        "nonce": "c0ae1c8e-c7e7-469f-b252-86e6a0e7387e", 
        "signatureValue": "BavEll0/I1zpYw8XNi1bgVg/sCneO4Jugez
         8RwDg/+MCRVpjOboDoe4SxxKjkCOvKiCHGDvc4krqi6Z1n0 
         UfqzxGfmatCuFibcC1wpsPRdW+gGsutPTLzvueMWmFhw 
         YmfIFpbBu95t501+rSLHIEuujM/+PXr9Cky6Ed+W3JT24=
    }
}    
```

## Security Considerations
The IoTeX DID method covers both people and devices and the following security issues are considered:

### Private Key Compromise
The ownership of DIDs is solely based on the knowledge of private keys. Hence the secure storage of private keys is critical. To minimize the risk of private key compromise, using security hardware for key storage is highly recommended. Furthermore, in the case that private key is lost or stolen, an identifier recovery mechanism needs to be in place for enabling a legitimate entity to regain control of the identifier. 

### Replay and Impersonation Attacks
A malicious service provider might collect credentials of people and devices with the purpose of impersonating other legitimate entities to gain access for certain services. This issue can be addressed by service providers conducting a challenge-response process with the request entities, thereby verifying that a DID is associated with the request entity.  

### Integrity of DID Documents and Credentials
The IoTeX DID method allows people and devices to store the DID documents and credentials in the storage of their choices. The integrity of this information is ensured by the IoTeX blockchain.

### Denial of Service Attacks
An attacker might launch the Denial of Service (DoS) attacks to prevent access of certain DID documents and/or credentials. To mitigate this issue, DID owners may deploy data redundancy countermeasures by storing multiple copies of DID documents and/or credentials in different storage locations.

### Smart Contract Flaws 
The IoTeX DIDs method is implemented by smart contracts on the IoTeX blockchain. Those contracts will go through stringent audits and tests to mitigate the potential security risks.

### Quantum Computer Threats
The private/public key pair used in IoTeX DID method is based on elliptic curve cryptosystems (ECC). The future occurrence of powerful quantum computers will render ECC insecure. Under such a circumstance, a key updating mechanism is needed to generate a new key pair based on the post-quantum cryptographic algorithm selected by the National Institute of Standards and Technology (NIST) for replacing the old one.

## Privacy Considerations
The IoTeX DID method involves privacy protection for people and devices and the following privacy issues are considered:

### Private Data Leakage
When a person shares his/her personal information with a service provider, the service provider might leak the data to third parties without users’ consent. This issue can be minimized by only sharing the necessary data during the interaction with the service provider rather than exposing the full credential. Advanced cryptographic techniques such as zero-knowledge proofs might be employed to address this issue under certain circumstances.  

### People/Device Tracking
Since IoTeX DIDs are managed by smart contracts and all the interactions with the smart contracts are recorded on the IoTeX blockchain. By analyzing the transactions on the blockchain, attackers might be able to review the interaction patterns among people or between people and devices. This issue can be minimized by constantly altering the DIDs of people and devices for different interactions according to a predefined policy. 

## Reference
1. Decentralized Identifiers (DIDs) v0.13 https://w3c-ccg.github.io/did-spec
2. Verifiable Claims https://www.w3.org/TR/verifiable-claims-data-model
3. RFC3552 https://tools.ietf.org/html/rfc3552
4. RFC6973 https://tools.ietf.org/html/rfc6973
