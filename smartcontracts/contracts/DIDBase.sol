pragma solidity >=0.4.22 <0.6.0;

import './IDID.sol';
import './DIDStorage.sol';
import './ownership/Ownable.sol';

// DIDBase is an abstract contract which implements IDID
contract DIDBase is IDID, Ownable {
    DIDStorage public db;
    string public prefix;

    constructor(address dbAddr) public {
        if (dbAddr == address(0x0)) {
            db = new DIDStorage();
        } else {
            db = DIDStorage(dbAddr);
        }
    }

    function transferDBOwnership(address newOwner) public onlyOwner {
        db.transferOwnership(newOwner);
    }

    function getHash(string memory did) public view returns (bytes32) {
        require(db.exist(did), "did does not exist");
        (,bytes32 h,) = db.get(did);

        return h;
    }

    function getURI(string memory did) public view returns (string memory) {
        require(db.exist(did), "did does not exist");
        (, , string memory uri) = db.get(did);
        return uri;
    }

    function getOwner(string memory did) public view returns (address) {
        require(db.exist(did), "did does not exist");
        (address owner, ,) = db.get(did);
        return owner;
    }
}
