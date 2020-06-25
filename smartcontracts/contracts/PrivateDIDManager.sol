pragma solidity >=0.4.22 <0.6.0;

import './DIDManagerBase.sol';
import './ownership/Whitelist.sol';

contract PrivateDIDManager is DIDManagerBase, Whitelist {
    constructor(bytes memory _prefix, address _dbAddr) DIDBase(_dbAddr, _prefix) public {
        addAddressToWhitelist(msg.sender);
    }

    function getInternalKey(bytes memory id) internal view returns (bytes20);

    function formDID(bytes memory id) public view returns (bytes memory);

    function createDID(bytes memory id, bytes32 h, bytes memory uri) public onlyWhitelisted {
        internalCreateDID(formDID(id), getInternalKey(id), msg.sender, h, uri);
    }

    function updateDID(bytes memory id, bytes32 h, bytes memory uri) public onlyWhitelisted {
        internalUpdateDID(formDID(id), getInternalKey(id), msg.sender, h, uri);
    }

    function deleteDID(bytes memory id) public onlyWhitelisted {
        internalDeleteDID(formDID(id), getInternalKey(id), msg.sender);
    }
}