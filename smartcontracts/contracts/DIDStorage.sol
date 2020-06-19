pragma solidity >=0.4.22 <0.6.0;

import './ownership/Ownable.sol';

contract DIDStorage is Ownable {
    struct DID {
        bool exist;
        address owner;
        bytes32 hash;
        string uri;
    }
    mapping(string => DID) dids;

    function exist(string memory did) public view returns (bool) {
        return dids[did].exist;
    }

    function upsert(string memory did, address owner, bytes32 hash, string memory uri) public onlyOwner {
        dids[did] = DID(true, owner, hash, uri);
    }

    function deactivate(string memory did) public onlyOwner {
        dids[did].exist = false;
    }

    function get(string memory did) public view returns (address, bytes32, string memory) {
        return (dids[did].owner, dids[did].hash, dids[did].uri);
    }
}
