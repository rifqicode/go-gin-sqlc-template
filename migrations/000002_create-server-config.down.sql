-- Drop the server_configs table and all associated objects
DROP TABLE IF EXISTS server_configs CASCADE;

-- Remove any custom types or functions if they were created (none in this case)

-- Note: The CASCADE option will automatically remove any dependent objects,
-- such as indexes, constraints, and triggers associated with the table.

