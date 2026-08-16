CREATE INDEX "events_session_path_definition_time" ON "public"."events" ("session_id", "path", "definition_id", "time" DESC);
CREATE INDEX "events_time_session_definition" ON "public"."events" ("time" DESC, "session_id", "definition_id");
CREATE INDEX "events_definition_time_path_session" ON "public"."events" ("definition_id", "time" DESC, "path", "session_id");
CREATE INDEX "sessions_site_enter_time" ON "public"."sessions" ("site_id", "enter_time");
CREATE INDEX "sessions_site_enter_day" ON "public"."sessions" ("site_id", "enter_day");
CREATE INDEX "sessions_site_enter_hour" ON "public"."sessions" ("site_id", "enter_hour");
