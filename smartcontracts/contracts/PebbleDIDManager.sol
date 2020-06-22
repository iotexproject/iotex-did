pragma solidity >=0.4.22 <0.6.0;

import './PrivateDIDManager.sol';

contract PebbleDIDManager is PrivateDIDManager {
    constructor(address _dbAddr) PrivateDIDManager("did:io:pebble:", _dbAddr) public {}

    function formDID(string memory id) public view returns (string memory) {
        // TODO: verify id format
        return string(abi.encodePacked(prefix, id));
    }
    /*
    function verifyDID(string memory did) public view returns (bool) {
        bytes memory ds = bytes(did);
        if (ds.length <= prefixBytes.length) {
            return false;
        }
        uint i = 0;
        // check prefix
        for (i = 0; i < prefixBytes.length; i++) {
            if (ds[i] != prefixBytes[i]) {
                return false;
            }
        }
        // check did format
        for (; i < ds.length; i++) {
            if (!(ds[i] >= 'a' && ds[i] <= 'z' || ds[i] >= 'A' && ds[i] <= 'Z' || ds[i] >= '0' && ds[i] <= '9')) {
                return false;
            }
        }
        return true;
    }
    */
}