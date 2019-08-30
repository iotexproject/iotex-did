pragma solidity ^0.4.24;

import './SelfManagedDID.sol';

contract IoTeXDID is SelfManagedDID {
    modifier onlyDIDOwner(string didInput) {
        string memory didString = generateDIDString();
        if (bytes(didInput).length > 0) {
            require(compareStrings(didInput, didString), "caller does not own the given did");
        }
        require(dids[didString].exist, "caller is not a did owner");
        _;
    }

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

    function updateHash(string did, bytes32 hash) public onlyDIDOwner(did) {
        dids[generateDIDString()].hash = hash;
    }

    function updateURI(string did, string uri) public onlyDIDOwner(did) {
        dids[generateDIDString()].uri = uri;
    }

    function deleteDID(string did) public onlyDIDOwner(did) {
        dids[generateDIDString()].exist = false;
    }

    function getHash(string did) public view returns (bytes32) {
        require(dids[did].exist, "did does not exist");
        return dids[did].hash;
    }

    function getURI(string did) public view returns (string) {
        require(dids[did].exist, "did does not exist");
        return dids[did].uri;
    }

    function generateDIDString() private view returns (string) {
        return string(abi.encodePacked(didPrefix, addrToString(msg.sender)));
    }

    function addrToString(address _addr) internal pure returns(string) {
        bytes32 value = bytes32(uint256(_addr));
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(43);
        str[0] = ' ';
        str[1] = '0';
        str[2] = 'x';
        for (uint i = 0; i < 20; i++) {
            str[3+i*2] = alphabet[uint(value[i + 12] >> 4)];
            str[4+i*2] = alphabet[uint(value[i + 12] & 0x0f)];
        }
        return string(str);
    }

    function compareStrings (string a, string b) internal pure returns (bool) {
        return (keccak256(abi.encodePacked((a))) == keccak256(abi.encodePacked((b))));
    }
}
