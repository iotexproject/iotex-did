pragma solidity ^0.4.24;

import './ownership/Whitelist.sol';
import './SelfManagedDID.sol';

contract DecentralizedIdentifier is Whitelist {
    mapping(bytes32 => address) public nameSpaceToSelfManagedAddress;

    event RegisterSelfManagedContract(bytes32 indexed nameSpace, address addr);
    event DeregisterSelfManagedContract(bytes32 indexed nameSpace);
    event UpdateSelfManagedContract(bytes32 indexed nameSpace, address addr);

    constructor() public {
        addAddressToWhitelist(msg.sender);
    }

    function registerSelfManagedContract(bytes32 nameSpace, address addr) public onlyWhitelisted {
        require(nameSpaceToSelfManagedAddress[nameSpace] == address(0), "name space is already registered");
        require(addr != address(0), "invalid contract address");
        nameSpaceToSelfManagedAddress[nameSpace] = addr;

        emit RegisterSelfManagedContract(nameSpace, addr);
    }

    function deregisterSelfManagedContract(bytes32 nameSpace) public onlyWhitelisted {
        require(nameSpaceToSelfManagedAddress[nameSpace] != address(0), "name space is not registered");
        nameSpaceToSelfManagedAddress[nameSpace] = address(0);

        emit DeregisterSelfManagedContract(nameSpace);
    }

    function updateSelfManagedContractAddress(bytes32 nameSpace, address addr) public onlyWhitelisted {
        require(nameSpaceToSelfManagedAddress[nameSpace] != address(0), "name space is not registered");
        require(addr != address(0), "invalid contract address");
        nameSpaceToSelfManagedAddress[nameSpace] = addr;

        emit UpdateSelfManagedContract(nameSpace, addr);
    }

    function getHash(bytes32 nameSpace, string did) public view returns (bytes32) {
        return getSelfManagedContract(nameSpace).getHash(did);
    }

    function getURI(bytes32 nameSpace, string did) public view returns (string) {
        return getSelfManagedContract(nameSpace).getURI(did);
    }

    function getSelfManagedContract(bytes32 nameSpace) private view returns (SelfManagedDID) {
        require(nameSpaceToSelfManagedAddress[nameSpace] != address(0), "name space does not exist");
        return SelfManagedDID(nameSpaceToSelfManagedAddress[nameSpace]);
    }
}