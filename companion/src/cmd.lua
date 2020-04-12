-- the entrypoint for chat based interactions
function JeevesAddon:ParseCmd(input)
    -- remove any slashes from the command
    input = string.trim(input, " ")

    -- /jeeves
    if input == "" or not input then
        return JeevesAddon:RootCmd()
    end

    -- /jeeves export
    if input == "export" then
        return JeevesAddon:ExportCmd()
    end

    if input == "reset" then
        return JeevesAddon:Reset()
    end

    -- we did not recognize the command
    print("Unrecognized command: \"" .. input .."\".  Please try again.")
end

-- a command with no inputs
function JeevesAddon:RootCmd()
    print(JeevesAddonTitle)
    print("|cFF80FF80/jeeves export|r - |cFFFF8080export the changes your inventories|r")
end

-- the command to export the delta between the last known inventory for this


function JeevesAddon:Reset()
    local existingBank = cachedBank()
    -- clear the cached bank if it exists
    if existingBank ~= nil and getn(existingBank) > 0 then
        resetCachedBank()
    end
    -- make sure we clear any export history
    LatestExports[UnitGUID("player")] = {}

    -- tell them we're done
    print("Successfully reset your bank data")
end
