pragma solidity >=0.4.22 <0.6.0;

contract Agentable {

    function getDeleteAuthMessage(bytes memory did, address agent) public view returns (bytes memory) {
        return abi.encodePacked(
            "I authorize ", addrToString(agent), " to delete DID ", did, " in contract ", addrToString(address(this)));
    }

    function getCreateAuthMessage(bytes memory did, bytes32 h, bytes memory uri, address agent) public view returns (bytes memory) {
        return abi.encodePacked(
            "I authorize ", addrToString(agent), " to create DID ", did, " in contract with ", addrToString(address(this)), " (", h, ", ", uri, ")");
    }

    function getUpdateAuthMessage(bytes memory did, bytes32 h, bytes memory uri, address agent) public view returns (bytes memory) {
        return abi.encodePacked(
            "I authorize ", addrToString(agent), " to update DID ", did, " in contract to ", addrToString(address(this)), " (", h, ", ", uri, ")");
    }

    function addrToString(address _addr) internal pure returns(string memory) {
        bytes32 value = bytes32(uint256(_addr));
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(42);
        str[0] = '0';
        str[1] = 'x';
        for (uint i = 0; i < 20; i++) {
            str[2+i*2] = alphabet[uint8(value[i + 12] >> 4)];
            str[3+i*2] = alphabet[uint8(value[i + 12] & 0x0f)];
        }
        return string(str);
    }

    function uint2str(uint i) internal pure returns (string memory) {
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
            b[k--] = byte(uint8(48 + i % 10));
            i /= 10;
        }
        return string(b);
    }


    function getSigner(bytes memory message, bytes memory signature) internal pure returns (address) {
        return recover(keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n", uint2str(message.length), message)), signature);
    }

    function recover(bytes32 hash, bytes memory signature)
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
}
