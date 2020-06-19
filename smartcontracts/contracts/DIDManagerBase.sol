pragma solidity >=0.4.22 <0.6.0;

import './DIDBase.sol';

// DIDManagerBase is an implementation of IDID which is only updatable by owner
contract DIDManagerBase is DIDBase {
    function internalCreateDID(string memory did, address didOwner, bytes32 h, string memory uri) internal {
        require(!db.exist(did), "duplicate did");
        db.upsert(did, didOwner, h, uri);

        emit DIDCreated(didOwner, did, h, uri);
    }

    function internalUpdateDID(string memory did, address didOwner, bytes32 h, string memory uri) internal {
        require(db.exist(did), "did does not exist");
        (address ownerAddr, ,) = db.get(did);
        require(ownerAddr == didOwner, "no permission");
        db.upsert(did, didOwner, h, uri);
        emit DIDUpdated(didOwner, did, h, uri);
    }

    function internalDeleteDID(string memory did, address didOwner) internal {
        require(db.exist(did), "did does not exist");
        (address ownerAddr, ,) = db.get(did);
        require(ownerAddr == didOwner, "no permission");
        db.deactivate(did);
        emit DIDDeleted(msg.sender, did);
    }
}