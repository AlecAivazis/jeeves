local AceGUI = LibStub("AceGUI-3.0")

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

    -- /jeeves register
    if input == "register" then
        return JeevesAddon:RegisterCmd()
    end

    -- /jeeves unregister
    if input == "unregister" then
        return JeevesAddon:UnRegisterCmd()
    end

    -- we did not recognize the command
    print("Unrecognized command: \"" .. input .."\".  Please try again.")
end

-- a command with no inputs
function JeevesAddon:RootCmd()
    print("Jeeves Companion v0.0.0")
    print("|cFF80FF80/jeeves register|r - |cFFFF8080register the current character as a bank alt|r")
    print("|cFF80FF80/jeeves unregister|r - |cFFFF8080unregister the current character as a bank alt|r")
    print("|cFF80FF80/jeeves export|r - |cFFFF8080export the changes your inventories|r")
end

-- a command that registers the current character as a bank alt
function JeevesAddon:RegisterCmd()
    -- if we haven't added a character before
    if CharacterInventories == nil then
        -- create the empty table of characters
        CharacterInventories = {}
    end

    -- add the players GUID to the table
    CharacterInventories[UnitGUID("player")] = {}

    print("Registered", GetUnitName("player"), "as a bank alt.")
end

-- a command that registers the current character as a bank alt
function JeevesAddon:UnRegisterCmd()
    -- if we haven't added a character before
    if CharacterInventories == nil then
        -- create the empty table of characters
        CharacterInventories = {}
    end

    -- if the player is a registered bank alt
    if IsBankAlt() then
        -- remove the entry for that user in the table
        CharacterInventories[UnitGUID("player")] = nil

        -- confirm our action with the user
        print("Unregistered", GetUnitName("player"), "as a bank alt.")
    else
        print(GetUnitName("player"), "is not a bank alt.")
    end
end

-- the command to export the delta between the last known inventory for this
-- bank and the current inventory
function JeevesAddon:ExportCmd()
    -- if we haven't added a character before
    if CharacterInventories == nil then
        -- create the empty table of characters
        CharacterInventories = {}
    end

    print("Exporting inventory...")

    -- build up the command the user needs to submit
    local command = "!deposit "
    for itemID, count in pairs(CachedBank()) do
        itemName = GetItemInfo(itemID)
        command = command ..     count ..    "x " ..    itemName ..    ","
    end
    command = command:sub(0, command:len()-1)

    -- we need to create a frame with the command
    local commandFrame = AceGUI:Create("Frame");
    commandFrame:SetWidth(500)
    commandFrame:SetHeight(125)
    commandFrame:SetTitle("Inventory Export")
    commandFrame:EnableResize(false)


    -- add the command in an edit box inside of the frame
    local editBox = AceGUI:Create("EditBox")
    editBox:SetWidth(450)
    editBox:SetHeight(50)
    editBox:SetText(command)
    commandFrame:AddChild(editBox)
end
