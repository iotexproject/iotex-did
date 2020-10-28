pragma solidity >=0.4.24;

import './SelfManagedDeviceDID.sol';
import './ownership/Whitelist.sol';

contract MockDeviceDID is SelfManagedDeviceDID, Whitelist {
    modifier onlyDeviceOwner(string uuid, bytes proof) {
        string memory didString = generateDIDString(uuid);
        require(dids[didString].exist, "did does not exist");
        verifyProof(uuid, proof);
        _;
    }

    string constant didPrefix = "did:io:mock:";
    struct DID {
        bool exist;
        bytes32 hash;
        string uri;
    }
    mapping(string => DID) dids;

    address public cloudServiceAddr;

    event CreateDID(address indexed owner, string indexed uuid, string didString);
    event UpdateHash(string indexed didString, bytes32 hash);
    event UpdateURI(string indexed didString, string uri);
    event DeleteDID(string indexed didString);

    constructor(address _cloudServiceAddr) public {
        addAddressToWhitelist(msg.sender);
        cloudServiceAddr = _cloudServiceAddr;
    }

    function createDID(string uuid, bytes proof, bytes32 hash, string uri) public {
        verifyProof(uuid, proof);
        string memory resultDID = generateDIDString(uuid);
        require(!dids[resultDID].exist, "did already exists");
        dids[resultDID] = DID(true, hash, uri);

        emit CreateDID(msg.sender, uuid, resultDID);
    }

    function updateHash(string uuid, bytes proof, bytes32 hash) public onlyDeviceOwner(uuid, proof) {
        string memory did = generateDIDString(uuid);
        dids[did].hash = hash;
        emit UpdateHash(did, hash);
    }

    function updateURI(string uuid, bytes proof, string uri) public onlyDeviceOwner(uuid, proof) {
        string memory did = generateDIDString(uuid);
        dids[did].uri = uri;
        emit UpdateURI(did, uri);
    }

    function deleteDID(string uuid, bytes proof) public onlyDeviceOwner(uuid, proof) {
        string memory did = generateDIDString(uuid);
        dids[did].exist = false;
        emit DeleteDID(did);
    }

    function setCloudServiceAddr(address _cloudServiceAddr) public onlyWhitelisted {
        cloudServiceAddr = _cloudServiceAddr;
    }

    function getHash(string did) public view returns (bytes32) {
        require(dids[did].exist, "did does not exist");
        return dids[did].hash;
    }

    function getURI(string did) public view returns (string) {
        require(dids[did].exist, "did does not exist");
        return dids[did].uri;
    }

    function verifyProof(string uuid, bytes proof) internal view {
         bytes memory message = abi.encodePacked("I authorize", addrToString(msg.sender), " to register ", uuid);
         require(recover(toEthPersonalSignedMessageHash(message), proof) == cloudServiceAddr, "invalid proof");
    }

    function recover(bytes32 hash, bytes signature)
        internal
        pure
        returns (address)
    {
        bytes32 r;
        bytes32 s;
        uint8 v;
        // Check the signature length
        if (signature.length != 65) {
            return (address(0));
        }
        // Divide the signature in r, s and v variables with inline assembly.
        assembly {
            r := mload(add(signature, 0x20))
            s := mload(add(signature, 0x40))
            v := byte(0, mload(add(signature, 0x60)))
        }
        // Version of signature should be 27 or 28, but 0 and 1 are also possible versions
        if (v < 27) {
            v += 27;
        }
        // If the version is correct return the signer address
        if (v != 27 && v != 28) {
            return (address(0));
        }
        return ecrecover(hash, v, r, s);
    }

    function generateDIDString(string uuid) private view returns (string) {
        return string(abi.encodePacked(didPrefix, uuid));
    }

    function uint2str(uint i) internal pure returns (string) {
        if (i == 0) {
            return "0";
        }
        uint j = i;
        uint length;
        while (j != 0){
            length++;
            j /= 10;
        }
        bytes memory b = new bytes(length);
        uint k = length - 1;
        while (i != 0){
            b[k--] = byte(48 + i % 10);
            i /= 10;
        }
        return string(b);
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

    function toEthPersonalSignedMessageHash(bytes _msg) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n", uint2str(_msg.length), _msg));
    }
}