Feature: event
		 In order to create event
		 As GRPC client of event service
		 I want to create event in event service

		 Scenario: create event
		        When I call grpc event method Create
		        Then The error should be nil
		        And The create response success should be true

		 Scenario: update event
		        When I call grpc event method Update
		        Then The error should be nil
		        And The update response success should be true