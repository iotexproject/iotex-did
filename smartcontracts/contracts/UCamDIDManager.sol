pragma solidity >=0.4.22 <0.6.0;

import './Agentable.sol';
import './DIDManagerBase.sol';

contract UCamDIDManager is Agentable, DIDManagerBase {

    constructor(address _dbAddr) DIDBase(_dbAddr) public {
        prefix = "ucam";
    }

    function formDID(string memory uid) public view returns (string memory) {
        bytes memory ds = bytes(uid);
        require(ds.length == 20, "invalid uid length");
        for (uint i = 0; i < ds.length; i++) {
            require(ds[i] >= 'A' && ds[i] <= 'Z' || ds[i] >= '0' && ds[i] <= '9', "invalid uid format");
        }
        return string(abi.encodePacked(prefix, uid));
    }

    function createDIDByAgent(string memory uid, bytes32 h, string memory uri, address authorizer, bytes memory auth) public {
        string memory did = formDID(uid);
        require(msg.sender == getSigner(getCreateAuthMessage(did, h, uri, msg.sender), auth), "invalid signature");
        internalCreateDID(did, authorizer, h, uri);
    }

    function updateDIDByAgent(string memory uid, bytes32 h, string memory uri, bytes memory auth) public {
        string memory did = formDID(uid);
        (address authorizer, ,) = db.get(did);
        require(msg.sender == getSigner(getUpdateAuthMessage(did, h, uri, msg.sender), auth), "invalid signature");
        internalUpdateDID(did, authorizer, h, uri);
    }

    function deleteDIDByAgent(string memory uid, bytes memory auth) public {
        string memory did = formDID(uid);
        (address authorizer, ,) = db.get(did);
        require(msg.sender == getSigner(getDeleteAuthMessage(did, msg.sender), auth), "invalid signature");
        internalDeleteDID(did, authorizer);
    }
}
