pragma solidity >=0.4.22 <0.6.0;

import './DIDBase.sol';

// DIDManagerBase is an implementation of IDID which is only updatable by owner
contract DIDManagerBase is DIDBase {
    function internalCreateDID(bytes memory did, bytes20 internalKey, address didOwner, bytes32 h, bytes memory uri) internal {
        require(!db.exist(internalKey), "duplicate DID");
        db.upsert(internalKey, true, didOwner, h, uri);

        emit DIDCreated(didOwner, string(did), h, string(uri));
    }

    function internalUpdateDID(bytes memory did, bytes20 internalKey, address didOwner, bytes32 h, bytes memory uri) internal {
        require(db.exist(internalKey), "DID does not exist");
        (address ownerAddr, ,) = db.get(internalKey);
        require(ownerAddr == didOwner, "no permission");
        db.upsert(internalKey, true, didOwner, h, uri);
        emit DIDUpdated(didOwner, string(did), h, string(uri));
    }

    function internalDeleteDID(bytes memory did, bytes20 internalKey, address didOwner) internal {
        require(db.exist(internalKey), "DID does not exist");
        (address ownerAddr, bytes32 h, bytes memory uri) = db.get(internalKey);
        require(ownerAddr == didOwner, "no permission");
        db.upsert(internalKey, false, didOwner, h, uri);
        emit DIDDeleted(msg.sender, string(did));
    }
}