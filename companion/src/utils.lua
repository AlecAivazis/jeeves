function cachedBank()
    return CharacterInventories[UnitGUID("player")]
end

function resetCachedBank()
    CharacterInventories[UnitGUID("player")] = {}
end

function latestExport()
    return LatestExports[UnitGUID("player")]
end

function getn (myTable)
    numItems = 0
    for k,v in pairs(myTable) do
        numItems = numItems + 1
    end

    return numItems
end

function shallowcopy(orig)
    local copy = {}

    for orig_key, orig_value in pairs(orig) do
        copy[orig_key] = orig_value
    end

    return copy
end
