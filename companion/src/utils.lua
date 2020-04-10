function IsBankAlt()
    return CachedBank() ~= nil
end

function CachedBank()
    -- if we haven't added a character before
    if CharacterInventories == nil then
        -- create the empty table of characters
        CharacterInventories = {}
    end

    return CharacterInventories[UnitGUID("player")]
end

function ResetCachedBank()
    -- if we haven't added a character before
    if CharacterInventories == nil then
        -- create the empty table of characters
        CharacterInventories = {}
    end

    CharacterInventories[UnitGUID("player")] = {}
end
