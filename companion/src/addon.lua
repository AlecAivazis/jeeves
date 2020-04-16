-- addon definition
JeevesAddon = LibStub("AceAddon-3.0"):NewAddon("Jeeves", "AceConsole-3.0")

-- invoked by ace when the addon is enabled
function JeevesAddon:OnEnable()
    -- register slash commands
    JeevesAddon:RegisterChatCommand("jeeves", "ParseCmd")
end

-- invoked by ace when the addon is disabled
function JeevesAddon:OnDisable()
    -- unregister slash commands
    JeevesAddon:UnregisterChatCommand("jeeves")
end
