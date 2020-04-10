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

    -- if we haven't added a character before
    if LatestExports == nil then
        -- create the empty table of characters
        LatestExports = {}
    end

    -- build up the command the user needs to submit
    local stem = "!deposit "


    -- we have to compute the total transactions to go from what we
    -- last had to what we have now
    local deposits, withdrawls = computeExports(LatestExports[UnitGUID("player")], CurrentInventory())

    -- we need to create a frame with the command
    local commandFrame = AceGUI:Create("Frame");
    commandFrame:SetWidth(500)
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
    local withdralMessage  = ""
    local depositMessage = ""
    local totalCommands = 0

    -- if there are deposits
    if deposits ~= nil then
        local commands, totalCount = buildCommands("!deposit ", deposits)

        -- if there is something to deposit
        if totalCount > 0 then
            totalCommands = totalCommands + table.getn(commands)

            depositMessage = "Depositing " .. totalCount .. " items."

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
    end

    -- if there are deposits
    if withdrawls ~= nil then
        local commands, totalCount = buildCommands("!withdraw ", withdrawls)

        -- if there is something to withdraw
        if totalCount > 0 then
            totalCommands = totalCommands + table.getn(commands)

            withdralMessage = "Withdrawing " .. totalCount .. " items."

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
    end

    -- if both of the messages are empty
    if withdralMessage == "" and depositMessage == "" then
        withdralMessage = "You have no items to export."
    else
        depositMessage = depositMessage .. " This must be done in "
                                        .. totalCommands .. " commands:"
    end

    -- use the messages
    text:SetText(withdralMessage .. depositMessage)

    -- save this as the latest export for the player
    LatestExports[UnitGUID("player")] = CurrentInventory()
end

function computeExports(latestExport, currentInventory)
    -- if we haven't seen anything before then its all deposits
    if latestExport == nil then
        return currentInventory, {}
    end

    -- we will return the list of deposits and withdrawls separately
    local deposits, withdrawls = {}, {}

    -- go over every entry in the current inventory
    for item, count in pairs(currentInventory) do
        -- look up this item in the latest export
        lastSeen = latestExport[item]

        -- if this is the first time we've seen the item
        if lastSeen == nil then
            -- use all of it
            deposits[item] = count

        -- if we have more than we last saw then deposit the extra
        elseif count > lastSeen then
            deposits[item] = count - lastSeen

        -- if we have less than we last saw what's missing has been withdrawn
        elseif count < lastSeen then
            withdrawls[item] = lastSeen - count
        end
    end

    -- we need to look for any items we saw last time we haven't seen now
    -- they were all withdrawn
    for item, count in pairs(latestExport) do
        if currentInventory[item] == nil then
            withdrawls[item] = count
        end
    end

    -- return both lists
    return deposits, withdrawls
end

function buildCommands(stem, entries)
    local commands = {}
    -- save a running count of the number of items we're exporting
    local totalCount = 0

    local currentCommand = stem
    for itemID, count in pairs(entries) do
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

    return commands, totalCount
end
