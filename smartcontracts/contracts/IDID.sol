pragma solidity >=0.4.22 <0.6.0;

interface DIDVerifier {
    function verify(string calldata did) external view returns (bool);
}

interface IDID {
    event DIDCreated(address indexed operator, string did, bytes32 hash, string uri);
    event DIDUpdated(address indexed operator, string indexed did, bytes32 hash, string uri);
    event DIDDeleted(address indexed operator, string indexed did);

    function getHash(bytes calldata) external view returns (bytes32);
    function getURI(bytes calldata) external view returns (bytes memory);
    function getOwner(bytes calldata) external view returns (address);
}