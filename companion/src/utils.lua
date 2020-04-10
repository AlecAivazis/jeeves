function IsBankAlt()
    return CachedBank() ~= nil
end

function CachedBank()
    return CharacterInventories[UnitGUID("player")]
end

function ResetCachedBank()
    CharacterInventories[UnitGUID("player")] = {}
end
