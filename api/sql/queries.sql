-- name: CreateAppointment :one
INSERT INTO appointments 
(patient_id, doctor_id, start_time, end_time, status) 
VALUES (
$1, $2, $3, $4, $5
) RETURNING id, patient_id, doctor_id, start_time, end_time, status;

-- name: GetAppointmentByID :one
SELECT * FROM appointments 
WHERE id = $1;

-- name: GetAppointmentsByDoctorID :many
SELECT * FROM appointments WHERE doctor_id = $1;

-- name: AddAppointmentCancellation :one
INSERT INTO appointment_cancellations
(appointment_id, patient_id, reason)
VALUES ($1, $2, $3)
RETURNING id, appointment_id, patient_id, reason;

-- name: UpdateAppointmentStatus :exec
UPDATE appointments SET status = $1 WHERE id = $2;

-- name: GetAppointmentsByDoctorIDAndStartAndEndTime :many
SELECT * FROM appointments 
WHERE doctor_id = $1 AND start_time = $2 AND end_time = $3;

-- name: UpdateAppointment :one
UPDATE appointments
SET patient_id = coalesce($1, patient_id), doctor_id = coalesce($2, doctor_id), start_time = coalesce($3, start_time), end_time = coalesce($4, end_time), status = coalesce($5, status)
WHERE id = $6
RETURNING id, patient_id, doctor_id, start_time, end_time, status;

-- name: GetDoctorSceduleforDayOfWeek :many
SELECT * FROM schedules WHERE doctor_id = $1 AND day_of_week = $2;

-- name: GetUpcomingAppointmentsByPatientID :many
SELECT * FROM appointments WHERE patient_id = $1 AND start_time > CURRENT_TIMESTAMP ORDER BY start_time DESC;

-- name: CreatePatient :one
INSERT INTO patients (name, email) VALUES ($1, $2) RETURNING id, name, email;

-- name: CreateDoctor :one
INSERT INTO doctors (name) VALUES ($1) RETURNING id, name;

-- name: CreateSchedule :one
INSERT INTO schedules (doctor_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id, doctor_id, day_of_week, start_time, end_time;


-- name: GetDoctors :many
SELECT * FROM doctors;

-- name: GetPatients :many
SELECT * FROM patients;

