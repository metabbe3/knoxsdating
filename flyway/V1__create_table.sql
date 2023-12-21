-- V1__create_user_entity_table.sql
CREATE TABLE IF NOT EXISTS "User" (
    "UserID" SERIAL PRIMARY KEY,
    "Username" VARCHAR(255) NOT NULL UNIQUE,
    "Email" VARCHAR(255) NOT NULL UNIQUE,
    "Password" VARCHAR(255) NOT NULL,
    "VerificationStatus" BOOLEAN DEFAULT false,
    "VerificationBadge" BOOLEAN DEFAULT false,
    "PremiumStatus" VARCHAR(50) DEFAULT 'Free' NOT NULL,
    "PremiumStartDate" TIMESTAMP,
    "PremiumEndDate" TIMESTAMP,
    "Gender" VARCHAR(20),
    "Company" VARCHAR(255),
    "School" VARCHAR(255),
    "JobTitle" VARCHAR(255),
    "VerifiedBadge" BOOLEAN DEFAULT false,
    "RedoCount" INT DEFAULT 1
);

CREATE TABLE IF NOT EXISTS "Profile" (
    "ProfileID" SERIAL PRIMARY KEY,
    "UserID" INT NOT NULL,
    "Photos" VARCHAR(255),
    "AboutMe" TEXT,
    "Interests" VARCHAR(255),
    "RelationshipGoals" TEXT,
    "Height" INT,
    "Language" VARCHAR(50),
    "ZodiacSign" VARCHAR(50),
    "EducationDetails" TEXT,
    "SocialMediaAccounts" VARCHAR(255),
    FOREIGN KEY ("UserID") REFERENCES "User"("UserID")
);

-- V3__create_swipe_history_v2_table.sql
CREATE TABLE IF NOT EXISTS "SwipeHistory" (
    "SwipeHistoryEntityID" SERIAL PRIMARY KEY,
    "SwiperUserID" INT NOT NULL,
    "SwipedUserID" INT NOT NULL,
    "SwipeDirection" VARCHAR(10) NOT NULL,
    "Timestamp" TIMESTAMP DEFAULT NOW() NOT NULL,
    "RedoCount" INT DEFAULT 0,
    "IsMatched" BOOLEAN DEFAULT false,
    FOREIGN KEY ("SwiperUserID") REFERENCES "User"("UserID"),
    FOREIGN KEY ("SwipedUserID") REFERENCES "User"("UserID")
);

-- V4__create_message_entity_table.sql
CREATE TABLE IF NOT EXISTS "Message" (
    "MessageID" SERIAL PRIMARY KEY,
    "SenderUserID" INT NOT NULL,
    "ReceiverUserID" INT NOT NULL,
    "MessageContent" TEXT NOT NULL,
    "Timestamp" TIMESTAMP NOT NULL,
    FOREIGN KEY ("SenderUserID") REFERENCES "User"("UserID"),
    FOREIGN KEY ("ReceiverUserID") REFERENCES "User"("UserID")
);

-- V5__create_report_entity_table.sql
CREATE TABLE IF NOT EXISTS "Report" (
    "ReportID" SERIAL PRIMARY KEY,
    "ReporterUserID" INT NOT NULL,
    "ReportedUserID" INT NOT NULL,
    "ReportContent" TEXT NOT NULL,
    "Timestamp" TIMESTAMP NOT NULL,
    FOREIGN KEY ("ReporterUserID") REFERENCES "User"("UserID"),
    FOREIGN KEY ("ReportedUserID") REFERENCES "User"("UserID")
);

-- V6__create_system_log_entity_table.sql
CREATE TABLE IF NOT EXISTS "SystemLog" (
    "LogID" SERIAL PRIMARY KEY,
    "LogType" VARCHAR(20) NOT NULL,
    "LogMessage" TEXT NOT NULL,
    "Timestamp" TIMESTAMP NOT NULL
);

-- V7__create_notification_entity_table.sql
CREATE TABLE IF NOT EXISTS "Notification" (
    "NotificationID" SERIAL PRIMARY KEY,
    "UserID" INT NOT NULL,
    "NotificationType" VARCHAR(50) NOT NULL,
    "Message" TEXT NOT NULL,
    "Timestamp" TIMESTAMP NOT NULL,
    "IsRead" BOOLEAN DEFAULT false,
    FOREIGN KEY ("UserID") REFERENCES "User"("UserID")
);

CREATE TABLE IF NOT EXISTS "Locationhistory" (
    "LocationID" SERIAL PRIMARY KEY,
    "UserID" INT NOT NULL,
    "Latitude" DOUBLE PRECISION,
    "Longitude" DOUBLE PRECISION,
    "Timestamp" TIMESTAMP NOT NULL,
    CONSTRAINT "FK_User_LocationHistory" FOREIGN KEY ("UserID") REFERENCES "User"("UserID")
);

CREATE TABLE IF NOT EXISTS "ProfileView" (
    "ProfileViewID" SERIAL PRIMARY KEY,
    "ViewerUserID" INT,
    "ShownUserID" INT,
    "Timestamp" TIMESTAMP,
    "DateOnly" DATE GENERATED ALWAYS AS (DATE("Timestamp")) STORED,
    CONSTRAINT "UniqueViewPerDay" UNIQUE ("ViewerUserID", "ShownUserID", "DateOnly")
);

-- V8__create_match_entity_table.sql
CREATE TABLE IF NOT EXISTS "Match" (
    "MatchID" SERIAL PRIMARY KEY,
    "UserID1" INT NOT NULL,
    "UserID2" INT NOT NULL,
    "Timestamp" TIMESTAMP NOT NULL,
    FOREIGN KEY ("UserID1") REFERENCES "User"("UserID"),
    FOREIGN KEY ("UserID2") REFERENCES "User"("UserID")
);


