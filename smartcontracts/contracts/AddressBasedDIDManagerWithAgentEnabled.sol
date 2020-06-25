pragma solidity >=0.4.22 <0.6.0;

import './AddressBasedDIDManager.sol';
import './Agentable.sol';

contract AddressBasedDIDManagerWithAgentEnabled is AddressBasedDIDManager, Agentable {

    constructor(bytes memory _prefix, address _dbAddr) AddressBasedDIDManager(_prefix, _dbAddr) public {}

    function registerByAgent(bytes32 h, bytes memory uri, address authorizer, bytes memory auth) public {
        bytes memory did = getDID(authorizer);
        bytes20 internalKey = bytes20(authorizer);
        require(!db.exist(internalKey), "duplicate DID");
        require(authorizer == getSigner(getCreateAuthMessage(did, h, uri, msg.sender), auth), "invalid signature");
        internalCreateDID(did, internalKey, authorizer, h, uri);
    }

    function updateByAgent(bytes memory did, bytes32 h, bytes memory uri, bytes memory auth) public {
        bytes20 internalKey = decodeInternalKey(did);
        require(db.exist(internalKey), "DID does not exist");
        (address authorizer, ,) = db.get(internalKey);
        require(authorizer == getSigner(getUpdateAuthMessage(did, h, uri, msg.sender), auth), "invalid signature");
        internalUpdateDID(did, internalKey, authorizer, h, uri);
    }

    function deregisterByAgent(bytes memory did, bytes memory auth) public {
        bytes20 internalKey = decodeInternalKey(did);
        require(db.exist(internalKey), "DID does not exist");
        (address authorizer, ,) = db.get(internalKey);
        require(authorizer == getSigner(getDeleteAuthMessage(did, msg.sender), auth), "invalid signature");
        internalDeleteDID(did, internalKey, authorizer);
    }
}
