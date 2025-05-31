CREATE TYPE task_status AS ENUM ('Pending', 'In Progress', 'Completed', 'Paused', 'Cancelled');
CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'urgent')
