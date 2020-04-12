-- need a frame to respond to events
local frame, events = CreateFrame("Frame"), {};

function events:BANKFRAME_OPENED(...)
    -- reset the cached inventory
    ResetCachedBank()

    -- export the bank container
    saveBag(BANK_CONTAINER, CachedBank())
    -- export every bank slot
    for slot = NUM_BAG_SLOTS+1, NUM_BAG_SLOTS+NUM_BANKBAGSLOTS do
        saveBag(slot, CachedBank())
    end
end

function saveBag(bagID, target)
    -- look up the number of slots in the bag
    for slot = 1, GetContainerNumSlots(bagID) do
        -- look up the item information at the slot
        _, itemCount, _, _, _, _, _, _, _, itemID = GetContainerItemInfo(bagID, slot)

        -- if we have an item in this slot
        if itemCount ~= nil then
            -- if we don't have an entry in the user's inventory for the
            if target[itemID] == nil then
                target[itemID] = itemCount
            else
                target[itemID] = target[itemID] + itemCount
            end
        end
    end
end

function CurrentInventory()
    -- lets build up a table of the players inventory
    local inventory = {}

    -- if there is a cached bank
    if CachedBank() ~= nil then
        for itemID, count in pairs(CachedBank()) do
            inventory[itemID] = count
        end
    end

    -- include every bag the character has
    for slot = 0, NUM_BAG_SLOTS do
        saveBag(slot, inventory)
    end

    -- we're done
    return inventory
end


-- implementation details


-- register the event handlers
frame:SetScript("OnEvent", function(self, event, ...)
 events[event](self, ...); -- call one of the functions above
end);
for k, v in pairs(events) do
 frame:RegisterEvent(k); -- Register all events for which handlers have been defined
end
