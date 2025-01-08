-- +goose Up
-- +goose StatementBegin
-- Create function to delete parent if no associated student
CREATE OR REPLACE FUNCTION delete_parent_with_no_associated_student() 
RETURNS TRIGGER AS $$
BEGIN
    -- Check if there is no other student that refers to the parent
    IF NOT EXISTS (
        SELECT 1 FROM students WHERE parent_uuid = OLD.parent_uuid AND deleted_at IS NULL
    ) THEN
        -- If no student refers to the parent, delete the parent
        UPDATE users
        SET deleted_at = NOW(), deleted_by = 'Auto Delete'
        WHERE user_uuid = OLD.parent_uuid;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to check and delete parent if no associated student
CREATE TRIGGER check_and_delete_parent_with_no_associated_student
AFTER UPDATE ON students
FOR EACH ROW
WHEN (OLD.deleted_at IS DISTINCT FROM NEW.deleted_at)
EXECUTE FUNCTION delete_parent_with_no_associated_student();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop trigger and function
DROP TRIGGER IF EXISTS check_and_delete_parent_with_no_associated_student ON students;
DROP FUNCTION IF EXISTS delete_parent_with_no_associated_student;
-- +goose StatementEnd