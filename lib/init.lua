-- golang-resource-pulsar-library main module.
-- Renders a Pulsar client and producer setup into the service's messaging package.
--
-- The calling archetype is responsible for adding the corresponding
-- Go module dependency:
--   github.com/apache/pulsar-client-go/pulsar
--
-- API (called from a parent archetype):
--   local pulsar = require("golang-resource-pulsar")
--   pulsar.render(context, { destination = context:get("project-name") })
--
-- Context contract (no required keys beyond what the calling archetype provides).

local M = {}

function M.render(context, opts)
    opts = opts or {}
    local d = opts.destination
    if d and d ~= "" then
        directory.render("contents", context, { destination = d })
    else
        directory.render("contents", context)
    end
    return context
end

return M
