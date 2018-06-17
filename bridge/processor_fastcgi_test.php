<?php

function accept()
{
    http_response_code(202);
    exit(0);
}

function fail($message)
{
    http_response_code(400);
    file_put_contents("php://stderr", $message);
    exit(0);
}

switch (@$_SERVER["TEST"]) {
    case "BODYSIZE10":
        $length = strlen(file_get_contents("php://input"));
        if ($length == 10) {
            accept();
        }

        fail(sprintf("Response size does not match: want %d, got %d\n", 10, $length));
        break;
    case "ENVVAR":
        if ($_SERVER["HTTP_FOO"] == "BAR") {
            accept();
        }

        fail(sprintf("Request should contain environment variable HTTP_FOO with value \"BAR\""));
        break;
    case "ACCEPT":
        accept();
        break;
	default:
		fail("Environment variable TEST is not set");
}
