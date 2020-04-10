-- need a frame to respond to events
local frame, events = CreateFrame("Frame"), {};

-- the ids of the bank slots

function events:BANKFRAME_OPENED(...)
    -- if the player is a registered bank alt then we need to track the bank
    -- contents
    if IsBankAlt() then
        -- reset the cached inventory
        ResetCachedBank()

        -- export the bank container
        saveBag(BANK_CONTAINER)
        -- export every bag slot
        for slot = 0, NUM_BAG_SLOTS + NUM_BANKBAGSLOTS do
            saveBag(slot)
        end
    end
end

function saveBag(bagID)
    -- look up the number of slots in the bag
    for slot = 1, GetContainerNumSlots(bagID) do
        -- look up the item information at the slot
        _, itemCount, _, _, _, _, _, _, _, itemID = GetContainerItemInfo(bagID, slot)

        -- if we have an item in this slot
        if itemCount ~= nil then
            -- if we don't have an entry in the user's inventory for the
            if CachedBank()[itemID] == nil then
                CachedBank()[itemID] = itemCount
            else
                CachedBank()[itemID] = CachedBank()[itemID] + itemCount
            end
        end
    end
end




-- implementation details


-- register the event handlers
frame:SetScript("OnEvent", function(self, event, ...)
 events[event](self, ...); -- call one of the functions above
end);
for k, v in pairs(events) do
 frame:RegisterEvent(k); -- Register all events for which handlers have been defined
end
