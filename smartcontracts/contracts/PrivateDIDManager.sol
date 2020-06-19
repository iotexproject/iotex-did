pragma solidity >=0.4.22 <0.6.0;

import './DIDManagerBase.sol';
import './ownership/Whitelist.sol';

contract PrivateDIDManager is DIDManagerBase, Whitelist {
    bytes private prefixBytes;

    constructor(string memory _prefix, address _dbAddr) DIDBase(_dbAddr) public {
        prefix = _prefix;
        prefixBytes = bytes(_prefix);
    }

    function formDID(string memory id) public view returns (string memory);

    function createDID(string memory id, bytes32 h, string memory uri) public onlyWhitelisted {
        internalCreateDID(formDID(id), msg.sender, h, uri);
    }

    function updateDID(string memory id, bytes32 h, string memory uri) public onlyWhitelisted {
        internalUpdateDID(formDID(id), msg.sender, h, uri);
    }

    function deleteDID(string memory id) public onlyWhitelisted {
        internalDeleteDID(formDID(id), msg.sender);
    }
}