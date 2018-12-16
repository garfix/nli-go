<?php
/**
 * Calls "nli answer" and echoes the result in JSON.
 */

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";

session_start();
$sessionId = session_id();

$command = __DIR__ . '/nli';
$configPath = __DIR__ . '/../resources/dbpedia/config-online.json';
$fullCommand = sprintf('%s -s %s -c %s "%s"', $command, $sessionId, $configPath, $query);

$process = exec($fullCommand, $output);

header('content-type: application/json');
echo implode("\n", $output);
