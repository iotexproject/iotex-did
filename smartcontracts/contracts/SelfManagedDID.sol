pragma solidity ^0.4.24;


interface SelfManagedDID {
    function createDID(string id, bytes32 hash, string uri) public;
    function deleteDID(string did) public;
    function updateHash(string did, bytes32 hash) public;
    function updateURI(string did, string uri) public;
    function getHash(string did) public view returns (bytes32);
    function getURI(string did) public view returns (string);
}