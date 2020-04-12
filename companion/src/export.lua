local AceGUI = LibStub("AceGUI-3.0")

local DiscordMessageLimit = 2000

-- ExportCmd gives the user the list of commands necessary to update Jeeves with the
-- current inventory of the bank alt
function JeevesAddon:ExportCmd()
    -- we have to compute the total transactions to go from what we
    -- last had to what we have now
    local deposits, withdrawls = computeExports(LatestExports[UnitGUID("player")], currentInventory())
    -- get the list of commands corresponding to each
    local depositCommands, totalDeposits = buildCommands("!deposit ", deposits)
    local withdrawlCommands, totalWithdrawls = buildCommands("!withdraw ", withdrawls)

    -- compute the total number of commands we're gonna have to run
    local totalCommands = getn(depositCommands) + getn(withdrawlCommands)

    if totalCommands == 0 then
        print("Sorry, there are no items to export.")
        return
    end

    -- we need to create a frame with the command
    local commandFrame = AceGUI:Create("Frame");
    commandFrame:SetWidth(500)
    commandFrame:SetHeight(100 * (totalCommands + 1) )
    commandFrame:SetTitle("Inventory Export")
    commandFrame:SetStatusText(JeevesAddonTitle)
    commandFrame:EnableResize(false)

    local spacer = AceGUI:Create("Label")
    spacer:SetText(" ")
    spacer:SetFontObject(GameFontHighlight)
    commandFrame:AddChild(spacer)

    -- add some text to the frame to tell the user what they are looking at
    local text = AceGUI:Create("Label")
    text:SetFullWidth(true)
    text:SetFontObject(GameFontHighlight)
    commandFrame:AddChild(text)

    local message = ""

    -- if nothing was exported, tell the user
    if totalDeposits == 0 and totalWithdrawls == 0 then
        message = "You have no items to export."

    -- notify them of what is going in an out
    else
        if totalDeposits > 0 then
            message = message .. "Depositing " .. totalDeposits .. " items."
        end

        if totalWithdrawls > 0 then
            message = message .. " Withdrawing " .. totalWithdrawls .. " items."
        end
    end

    --  if it taks more than one command, tell them
    if totalCommands > 1 then
        message = message .. " This will take " .. totalCommands .. " separate commands."
    end

    -- update the message
    text:SetText(message)

    for _, command in pairs(depositCommands) do
        local editBox = AceGUI:Create("EditBox")
        editBox:SetWidth(450)
        editBox:SetHeight(50)
        editBox:SetText(command)
        commandFrame:AddChild(editBox)
    end
    for _, command in pairs(withdrawlCommands) do
        local editBox = AceGUI:Create("EditBox")
        editBox:SetWidth(450)
        editBox:SetHeight(50)
        editBox:SetText(command)
        commandFrame:AddChild(editBox)
    end

    -- add the button to confirm the export
    local button = AceGUI:Create("Button")
    button:SetText("Confirm")
    button:SetCallback("OnClick", function()
        updateLatestExport()
        commandFrame:Hide()
    end)
    commandFrame:AddChild(button)
end

-- computeExports computes the operations that update the bank (assumed to be at the state
-- specified by the first arg) to the state passed as the second arg
function computeExports(latestExport, currentInventory)
    -- if we haven't seen anything before then its all deposits
    if latestExport == nil or getn(latestExport) == 0 then
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

-- buildCommands builds the list of string commands to send to jeeves, assuminig
-- that a given command cannot exceed {DiscordMessageLimit}
function buildCommands(stem, entries)
    local commands = {}
    -- save a running count of the number of items we're exporting
    local totalCount = 0

    if entries == nil then
        return {}, 0
    end

    local currentCommand = stem
    for itemID, count in pairs(entries) do
        local itemName = GetItemInfo(itemID)
        -- the entry we are going to add for this item
        local depositEntry = count ..    "x " .. itemName .. ","
        -- increment the total
        totalCount = totalCount + count

        -- if this entry will bring us above the limit
        if (currentCommand:len() + depositEntry:len()) > DiscordMessageLimit then
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

    -- if there are no commands, return an empty list
    if totalCount == 0 then
        return {}, 0
    end

    return commands, totalCount
end

-- updateLatestExport saves the current inventory as the latest snapshot. This
-- should be used when the user has confirmed they sent the messages to Jeeves
function updateLatestExport()
    -- save this as the latest export for the player
    LatestExports[UnitGUID("player")] = {}
    for key, value in pairs(currentInventory()) do
        LatestExports[UnitGUID("player")][key] = value
    end
end
