pragma solidity >=0.4.22 <0.6.0;

import './IDID.sol';
import './DIDStorage.sol';
import './ownership/Ownable.sol';

// DIDBase is an abstract contract which implements IDID
contract DIDBase is IDID, Ownable {
    DIDStorage public db;

    constructor(address dbAddr, bytes memory prefix) public {
        if (dbAddr == address(0x0)) {
            db = new DIDStorage(prefix);
        } else {
            db = DIDStorage(dbAddr);
        }
    }

    function transferDBOwnership(address newOwner) public onlyOwner {
        db.transferOwnership(newOwner);
    }

    function hasPrefix(bytes memory _bytes, bytes memory _prefix) internal pure returns (bool) {
        if (_bytes.length <= _prefix.length) {
            return false;
        }
        for (uint i = 0; i < _prefix.length; i++) {
            if (_prefix[i] != _bytes[i]) {
                return false;
            }
        }
        return true;
    }

    function slice(bytes memory _bytes, uint256 _start) internal pure returns (bytes memory) {
        require(_bytes.length >= _start, "Read out of bounds");
        uint256 _length = _bytes.length - _start;
        bytes memory tempBytes;

        assembly {
            switch iszero(_length)
            case 0 {
                // Get a location of some free memory and store it in tempBytes as
                // Solidity does for memory variables.
                tempBytes := mload(0x40)

                // The first word of the slice result is potentially a partial
                // word read from the original array. To read it, we calculate
                // the length of that partial word and start copying that many
                // bytes into the array. The first word we copy will start with
                // data we don't care about, but the last `lengthmod` bytes will
                // land at the beginning of the contents of the new array. When
                // we're done copying, we overwrite the full first word with
                // the actual length of the slice.
                let lengthmod := and(_length, 31)

                // The multiplication in the next line is necessary
                // because when slicing multiples of 32 bytes (lengthmod == 0)
                // the following copy loop was copying the origin's length
                // and then ending prematurely not copying everything it should.
                let mc := add(add(tempBytes, lengthmod), mul(0x20, iszero(lengthmod)))
                let end := add(mc, _length)

                for {
                    // The multiplication in the next line has the same exact purpose
                    // as the one above.
                    let cc := add(add(add(_bytes, lengthmod), mul(0x20, iszero(lengthmod))), _start)
                } lt(mc, end) {
                    mc := add(mc, 0x20)
                    cc := add(cc, 0x20)
                } {
                    mstore(mc, mload(cc))
                }

                mstore(tempBytes, _length)

                //update free-memory pointer
                //allocating the array padded to 32 bytes like the compiler does now
                mstore(0x40, and(add(mc, 31), not(31)))
            }
            //if we want a zero-length slice let's just return a zero-length array
            default {
                tempBytes := mload(0x40)

                mstore(0x40, add(tempBytes, 0x20))
            }
        }

        return tempBytes;
    }

    function decodeInternalKey(bytes memory did) public view returns (bytes20);

    function getHash(bytes memory did) public view returns (bytes32) {
        bytes20 internalKey = decodeInternalKey(did);
        require(db.exist(internalKey), "DID does not exist");
        (,bytes32 h,) = db.get(internalKey);

        return h;
    }

    function getURI(bytes memory did) public view returns (bytes memory) {
        bytes20 internalKey = decodeInternalKey(did);
        require(db.exist(internalKey), "DID does not exist");
        (, , bytes memory uri) = db.get(internalKey);
        return uri;
    }

    function getOwner(bytes memory did) public view returns (address) {
        bytes20 internalKey = decodeInternalKey(did);
        require(db.exist(internalKey), "DID does not exist");
        (address owner, ,) = db.get(internalKey);
        return owner;
    }
}
