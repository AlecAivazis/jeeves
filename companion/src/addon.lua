-- addon definition
JeevesAddon = LibStub("AceAddon-3.0"):NewAddon("Jeeves", "AceConsole-3.0")

JeevesAddonTitle = "Jeeves Companion v0.0.0"

-- invoked by ace when the addon is enabled
function JeevesAddon:OnEnable()
    -- register slash commands
    JeevesAddon:RegisterChatCommand("jeeves", "ParseCmd")
end

function JeevesAddon:OnInitialize()
    -- if we haven't added a character before
    if CharacterInventories == nil then
        -- create the empty table of characters
        CharacterInventories = {}
    end

    -- if we haven't added a character before
    if LatestExports == nil then
        -- create the empty table of characters
        LatestExports = {}
    end

end

-- invoked by ace when the addon is disabled
function JeevesAddon:OnDisable()
    -- unregister slash commands
    JeevesAddon:ResetChatCommand("jeeves")
end
