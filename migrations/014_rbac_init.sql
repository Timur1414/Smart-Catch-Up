create role aggregator_service_role;
create role api_gateway_service_role;
create role auth_service_role;
create role generator_service_role;

grant usage on schema public to aggregator_service_role, api_gateway_service_role, auth_service_role, generator_service_role;

grant insert on block, block_part to aggregator_service_role;
grant update on notification, digest to aggregator_service_role;
grant select on notification, "user", settings, block, digest to aggregator_service_role;

grant select on block, block_part, digest, notification, notification_action, pg_enum, "user", settings to api_gateway_service_role;
grant update on "user", settings to api_gateway_service_role;

grant insert on "user", settings, digest to auth_service_role;
grant select on "user" to auth_service_role;
grant update on "user" to auth_service_role;

grant insert on notification, notification_action to generator_service_role;
grant select on "user", settings, notification to generator_service_role;

---- create above / drop below ----

revoke insert on block, block_part from aggregator_service_role;
revoke update on notification, digest from aggregator_service_role;
revoke select on notification, "user", settings from aggregator_service_role;

revoke select on block, block_part, digest, notification, notification_action, pg_enum, "user", settings from api_gateway_service_role;
revoke update on "user", settings from api_gateway_service_role;

revoke insert on "user", settings, digest from auth_service_role;
revoke select on "user" from auth_service_role;
revoke update on "user" from auth_service_role;

revoke insert on notification, notification_action from generator_service_role;
revoke select on "user", settings from generator_service_role;

revoke usage on schema public from aggregator_service_role, api_gateway_service_role, auth_service_role, generator_service_role;
