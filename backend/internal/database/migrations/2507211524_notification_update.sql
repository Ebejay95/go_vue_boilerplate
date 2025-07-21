-- internal/database/migrations/2507211524_notification_update.sql
-- Add enhanced notification features with user support

-- Add new columns to notifications table
ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
ADD COLUMN IF NOT EXISTS read BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS persistent BOOLEAN DEFAULT TRUE;

-- Update existing notifications to have default values
UPDATE notifications SET persistent = TRUE WHERE persistent IS NULL;
UPDATE notifications SET read = FALSE WHERE read IS NULL;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read);
CREATE INDEX IF NOT EXISTS idx_notifications_persistent ON notifications(persistent);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, read);

-- Create utility function to clean up old notifications
CREATE OR REPLACE FUNCTION cleanup_old_notifications(days_old INTEGER DEFAULT 30)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM notifications
    WHERE read = TRUE
    AND persistent = TRUE
    AND created_at < CURRENT_TIMESTAMP - (days_old || ' days')::INTERVAL;

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create function to get notification statistics
CREATE OR REPLACE FUNCTION get_notification_stats(p_user_id INTEGER)
RETURNS TABLE(
    total_count INTEGER,
    unread_count INTEGER,
    read_count INTEGER,
    notification_types JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::INTEGER as total_count,
        COUNT(CASE WHEN read = FALSE THEN 1 END)::INTEGER as unread_count,
        COUNT(CASE WHEN read = TRUE THEN 1 END)::INTEGER as read_count,
        COALESCE(
            jsonb_object_agg(type, type_count),
            '{}'::jsonb
        ) as notification_types
    FROM (
        SELECT
            type,
            COUNT(*) as type_count
        FROM notifications
        WHERE (p_user_id IS NULL OR user_id = p_user_id OR (p_user_id IS NOT NULL AND user_id IS NULL))
        AND persistent = TRUE
        GROUP BY type
    ) type_counts;
END;
$$ LANGUAGE plpgsql;

-- Add some user-specific notifications
INSERT INTO notifications (message, type, user_id, read, persistent) VALUES
    ('Your profile has been updated', 'success', 1, false, true),
    ('New feature available', 'info', 2, false, true),
    ('Account verified successfully', 'success', 2, true, true),
    ('System maintenance tonight', 'warning', NULL, false, true)
ON CONFLICT DO NOTHING;
