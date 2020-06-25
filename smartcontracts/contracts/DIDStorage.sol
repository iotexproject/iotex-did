pragma solidity >=0.4.22 <0.6.0;

import './ownership/Ownable.sol';

contract DIDStorage is Ownable {
    bytes private prefix;

    struct DID {
        bool exist;
        address owner;
        bytes32 hash;
        bytes uri;
    }
    mapping(bytes20 => DID) dids;

    constructor(bytes memory _prefix) public {
        setPrefix(_prefix);
    }

    function exist(bytes20 internalDID) public view returns (bool) {
        return dids[internalDID].exist;
    }

    function setPrefix(bytes memory _prefix) public onlyOwner {
        prefix = _prefix;
    }

    function getPrefix() public view returns (bytes memory) {
        return prefix;
    }

    function upsert(bytes20 internalDID, bool isExist, address owner, bytes32 hash, bytes memory uri) public onlyOwner {
        dids[internalDID] = DID(isExist, owner, hash, uri);
    }

    function get(bytes20 internalDID) public view returns (address, bytes32, bytes memory) {
        require(exist(internalDID), "DID does not exist");
        return (dids[internalDID].owner, dids[internalDID].hash, dids[internalDID].uri);
    }
}
