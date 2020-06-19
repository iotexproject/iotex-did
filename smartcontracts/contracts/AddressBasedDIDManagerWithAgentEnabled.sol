pragma solidity >=0.4.22 <0.6.0;

import './AddressBasedDIDManager.sol';
import './Agentable.sol';

contract AddressBasedDIDManagerWithAgentEnabled is AddressBasedDIDManager, Agentable {

    constructor(string memory _prefix, address _dbAddr) AddressBasedDIDManager(_prefix, _dbAddr) public {}

    function registerByAgent(bytes32 h, string memory uri, address authorizer, bytes memory auth) public {
        require(msg.sender == getSigner(getCreateAuthMessage(getDID(authorizer), h, uri, msg.sender), auth), "invalid signature");
        internalCreateDID(getDID(authorizer), authorizer, h, uri);
    }

    function updateByAgent(string memory did, bytes32 h, string memory uri, bytes memory auth) public {
        (address authorizer, ,) = db.get(did);
        require(msg.sender == getSigner(getUpdateAuthMessage(getDID(authorizer), h, uri, msg.sender), auth), "invalid signature");
        internalUpdateDID(did, authorizer, h, uri);
    }

    function deregisterByAgent(string memory did, bytes memory auth) public {
        (address authorizer, ,) = db.get(did);
        require(msg.sender == getSigner(getDeleteAuthMessage(getDID(authorizer), msg.sender), auth), "invalid signature");
        internalDeleteDID(did, authorizer);
    }
}
