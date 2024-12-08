-- +goose Up
-- +goose StatementBegin
-- Membuat fungsi untuk menghapus orang tua jika tidak ada siswa yang merujuk
CREATE OR REPLACE FUNCTION delete_parent_with_no_associated_student() 
RETURNS TRIGGER AS $$
BEGIN
    -- Cek apakah orang tua (parent_uuid) masih dirujuk oleh siswa lain
    IF NOT EXISTS (
        SELECT 1 FROM students WHERE parent_uuid = OLD.parent_uuid
    ) THEN
        -- Jika tidak ada siswa lain yang merujuk, hapus orang tua
        DELETE FROM users WHERE user_uuid = OLD.parent_uuid;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Membuat trigger pada tabel students untuk memanggil fungsi
CREATE TRIGGER check_and_delete_parent_with_no_associated_student
AFTER DELETE ON students
FOR EACH ROW
EXECUTE FUNCTION delete_parent_with_no_associated_student();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Menghapus trigger dan fungsi jika rollback
DROP TRIGGER IF EXISTS check_and_delete_parent_with_no_associated_student ON students;
DROP FUNCTION IF EXISTS delete_parent_with_no_associated_student;
-- +goose StatementEnd