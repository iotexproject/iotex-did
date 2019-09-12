# IoTeX DID Method Specification
### Second Draft 10 Sep 2019

A DID is an identifier which enables users to protect and control their identities on IoTeX blockchain. In IoTeX DID system, we allow each manufacture or namespace to store and manage DIDs through its self-managed DID smart contract. Every self-managed smart contract has to implement the SelfManagedDID interface defined as follows:

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
2.  Append the UUID (IoTeX account address) to the prefix. The UUID must match caller's address. If UUID is missing, the contract will use caller's address as UUID. 
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

## Deactivate a DID

A registered DID can be deactivated anytime as long as the caller's address matches the address within the DID string. Once a DID is deactivated, the metadata of the corresponding document cannot be updated. 

## DID Documentation Resolution

### Update

Whenever a DID document is modified, the document hash will change and need to be updated in the DID contract. If a user moves the document to a new storage location, the document URI will change and need to be updated in the DID contract. Every device can only update its own document associated with the DID created by itself.

### Query
Given a unique did string, everyone can query the DID document hash and URI whether it owns the DID or not.  

TBD
