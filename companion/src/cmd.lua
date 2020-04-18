local AceGUI = LibStub("AceGUI-3.0")

local MessageLimit = 2000

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

    -- we did not recognize the command
    print("Unrecognized command: \"" .. input .."\".  Please try again.")
end

-- a command with no inputs
function JeevesAddon:RootCmd()
    print("Jeeves Companion v0.0.0")
    print("|cFF80FF80/jeeves export|r - |cFFFF8080export the changes your inventories|r")
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
    local stem = "!deposit "

    -- we need to break up the commands in 2000 character messages
    local commands = {}
    -- save a running count of the number of items we're exporting
    local totalCount = 0

    local currentCommand = stem
    for itemID, count in pairs(CurrentInventory()) do
        -- the entry we are going to add for this item
        local depositEntry = count ..    "x " .. GetItemInfo(itemID) .. ","
        -- increment the total
        totalCount = totalCount + count

        -- if this entry will bring us above the limit
        if currentCommand:len() + depositEntry:len() > MessageLimit then
            -- add the current command to the list
            table.insert(commands, currentCommand:sub(0, currentCommand:len()-1))

            -- reset the current command
            currentCommand = stem .. depositEntry
        else
            -- add the entry to the running command
            currentCommand = currentCommand .. depositEntry
        end
    end

    -- add whatever command we were building up at the end
    table.insert(commands, currentCommand:sub(0, currentCommand:len()-1))

    -- we need to create a frame with the command
    local commandFrame = AceGUI:Create("Frame");
    commandFrame:SetWidth(500)
    commandFrame:SetHeight(100 * table.getn(commands))
    commandFrame:SetTitle("Inventory Export")
    commandFrame:EnableResize(false)

    local spacer = AceGUI:Create("Label")
    spacer:SetText(" ")
    spacer:SetFontObject(GameFontHighlight)
    commandFrame:AddChild(spacer)

    -- add some text to the frame to tell the user what they are looking at
    local text  = AceGUI:Create("Label")
    text:SetFullWidth(true)
    text:SetFontObject(GameFontHighlight)
    commandFrame:AddChild(text)

    -- the message to show
    local message  = "Exporting " .. totalCount .. " items. "
    -- if there is more than one message to show, explain it
    if table.getn(commands) > 1 then
        message = message .. "This must be done in " .. table.getn(commands) .. " messages."
    end

    -- use the messages
    text:SetText(message)

    -- we need an edit box for every command the banker needs
    -- to submit
    for _, command in pairs(commands) do
        local editBox = AceGUI:Create("EditBox")
        editBox:SetWidth(450)
        editBox:SetHeight(50)
        editBox:SetText(command)
        commandFrame:AddChild(editBox)
    end
end
