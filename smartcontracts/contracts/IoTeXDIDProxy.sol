pragma solidity >=0.4.21 <0.6.0;

import "./IoTeXDIDStorage.sol";
import "./ownership/Ownable.sol";

contract IoTeXDIDProxy is IoTeXDIDStorage,Ownable {
    // version=>contract address
    mapping(string => address) public allVersions;
    // version list
    string[] public versionList;
    // current version
    string public currentVersion;
    event Upgrade(string version, address addr);

    constructor(address addr) public {
        upgrade("0.0.1", addr);
    }

    function currentDIDAddress() public view returns (address) {
        return allVersions[currentVersion];
    }

    function upgrade(string memory version, address addr) public onlyOwner {
        require(currentDIDAddress() != addr && addr != address(0), "Old address is not allowed and contract address should not be 0x");
        require(isContract(addr), "Cannot set a logic contract to a non-contract address");
        require(bytes(version).length > 0, "Version should not be empty string");
        currentVersion = version;
        allVersions[currentVersion] = addr;
        versionList.push(currentVersion);
        emit Upgrade(currentVersion, addr);
    }

    function getDIDAddress(string memory version) public view returns(address) {
        require(bytes(version).length > 0, "Version should not be empty string");
        return allVersions[version];
    }

    function isContract(address account) internal view returns (bool) {
        // According to EIP-1052, 0x0 is the value returned for not-yet created accounts
        // and 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470 is returned
        // for accounts without code, i.e. `keccak256('')`
        bytes32 codehash;
        bytes32 accountHash = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470;
        // solhint-disable-next-line no-inline-assembly
        assembly { codehash := extcodehash(account) }
        return (codehash != accountHash && codehash != 0x0);
    }

    function() payable external {
        address _impl = currentDIDAddress();
        require(_impl != address(0), "did contract address not set");

        assembly {
            let ptr := mload(0x40)
            calldatacopy(ptr, 0, calldatasize)
            let result := delegatecall(gas, _impl, ptr, calldatasize, 0, 0)
            let size := returndatasize
            returndatacopy(ptr, 0, size)

            switch result
            case 0 { revert(ptr, size) }
            default { return(ptr, size) }
        }
    }
}
