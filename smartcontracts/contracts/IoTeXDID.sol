pragma solidity ^0.4.24;

contract IoTeXDID {
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

    event CreateDID(string indexed id, string didString);
    event UpdateHash(string indexed didString, bytes32 hash);
    event UpdateURI(string indexed didString, string uri);
    event DeleteDID(string indexed didString);

    function createDID(string id, bytes32 hash, string uri) public {
        if (bytes(id).length > 0) {
            require(compareStrings(id, addrToString(msg.sender)), "id does not match creator");
        }
        string memory resultDID = generateDIDString();
        require(!dids[resultDID].exist, "did already exists");
        dids[resultDID] = DID(true, hash, uri);

        emit CreateDID(toLower(addrToString(msg.sender)), resultDID);
    }

    function updateHash(string did, bytes32 hash) public onlyDIDOwner(did) {
        dids[generateDIDString()].hash = hash;
        emit UpdateHash(generateDIDString(), hash);
    }

    function updateURI(string did, string uri) public onlyDIDOwner(did) {
        dids[generateDIDString()].uri = uri;
        emit UpdateURI(generateDIDString(), uri);
    }

    function deleteDID(string did) public onlyDIDOwner(did) {
        dids[generateDIDString()].exist = false;
        emit DeleteDID(generateDIDString());
    }

    function getHash(string did) public view returns (bytes32) {
        string memory didString = toLower(did);
        require(dids[didString].exist, "did does not exist");
        return dids[didString].hash;
    }

    function getURI(string did) public view returns (string) {
        string memory didString = toLower(did);
        require(dids[didString].exist, "did does not exist");
        return dids[didString].uri;
    }

    function generateDIDString() private view returns (string) {
        return string(abi.encodePacked(didPrefix, addrToString(msg.sender)));
    }

    function addrToString(address _addr) internal pure returns(string) {
        bytes32 value = bytes32(uint256(_addr));
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(42);
        str[0] = '0';
        str[1] = 'x';
        for (uint i = 0; i < 20; i++) {
            str[2+i*2] = alphabet[uint(value[i + 12] >> 4)];
            str[3+i*2] = alphabet[uint(value[i + 12] & 0x0f)];
        }
        return string(str);
    }

    function compareStrings (string a, string b) internal pure returns (bool) {
        return (keccak256(abi.encodePacked((toLower(a)))) == keccak256(abi.encodePacked((toLower(b)))));
    }

    function toLower(string str) internal pure returns (string) {
		bytes memory bStr = bytes(str);
		bytes memory bLower = new bytes(bStr.length);
		for (uint i = 0; i < bStr.length; i++) {
			if ((bStr[i] >= 65) && (bStr[i] <= 90)) {
				bLower[i] = bytes1(int(bStr[i]) + 32);
			} else {
				bLower[i] = bStr[i];
			}
		}
		return string(bLower);
	}
}
