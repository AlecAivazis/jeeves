-- need a frame to respond to events
local frame = CreateFrame("FRAME")

-- we need to track the inventory in every bank character's inventory
frame:RegisterEvent("BANKFRAME_OPENED")