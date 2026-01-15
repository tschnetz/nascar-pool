// 2026 NASCAR Cup Series Schedule with track, TV, and time info
export interface ScheduleEntry {
	raceNumber: number | null;
	date: string;
	name: string;
	track: string;
	stageLaps: string;
	startTime: string;
	tv: string;
	isSpecial: boolean;
	isOffWeek: boolean;
	isNonPoints: boolean;
}

export const schedule: ScheduleEntry[] = [
	{ raceNumber: null, date: 'Sun., Feb 1', name: 'Cook Out Clash', track: 'Bowman Gray Stadium', stageLaps: 'TBA', startTime: '8:00pm', tv: 'FOX', isSpecial: false, isOffWeek: false, isNonPoints: true },
	{ raceNumber: null, date: 'Thu., Feb 12', name: 'Daytona 500 Duels', track: 'Daytona International Speedway', stageLaps: 'N/A', startTime: '7:00pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: true },
	{ raceNumber: 1, date: 'Sun., Feb 15', name: 'DAYTONA 500', track: 'Daytona International Speedway', stageLaps: '65/130/200', startTime: '2:30pm', tv: 'FOX', isSpecial: true, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 2, date: 'Sun., Feb 22', name: 'Ambetter Health 400', track: 'EchoPark Speedway', stageLaps: '60/160/260', startTime: '3:00pm', tv: 'FOX', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 3, date: 'Sun., Mar 1', name: 'DuraMAX Grand Prix', track: 'Circuit of the Americas', stageLaps: '20/45/95', startTime: '3:30pm', tv: 'FOX', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 4, date: 'Sun., Mar 8', name: 'Straight Talk Wireless 500', track: 'Phoenix Raceway', stageLaps: '60/185/312', startTime: '3:30pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 5, date: 'Sun., Mar 15', name: 'Pennzoil 400', track: 'Las Vegas Motor Speedway', stageLaps: '80/165/267', startTime: '4:00pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 6, date: 'Sun., Mar 22', name: 'Goodyear 400', track: 'Darlington Raceway', stageLaps: '90/185/293', startTime: '3:00pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 7, date: 'Sun., Mar 29', name: 'Cook Out 400', track: 'Martinsville Speedway', stageLaps: '80/180/400', startTime: '3:30pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: null, date: 'Sun., Apr 5', name: 'Easter Break', track: 'OFF WEEK', stageLaps: '', startTime: '', tv: '', isSpecial: false, isOffWeek: true, isNonPoints: false },
	{ raceNumber: 8, date: 'Sun., Apr 12', name: 'Food City 500', track: 'Bristol Motor Speedway', stageLaps: '125/250/500', startTime: '3:00pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 9, date: 'Sun., Apr 19', name: 'AdventHealth 400', track: 'Kansas Speedway', stageLaps: '80/165/267', startTime: '2:00pm', tv: 'FOX', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 10, date: 'Sun., Apr 26', name: "Jack Link's 500", track: 'Talladega Superspeedway', stageLaps: '60/120/188', startTime: '3:00pm', tv: 'FOX', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 11, date: 'Sun., May 3', name: 'Würth 400', track: 'Texas Motor Speedway', stageLaps: '80/165/267', startTime: '3:30pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 12, date: 'Sun., May 10', name: 'Go Bowling at the Glen', track: 'Watkins Glen International', stageLaps: '20/40/90', startTime: '3:00pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: null, date: 'Sun., May 17', name: 'All-Star Race', track: 'Dover Motor Speedway', stageLaps: 'TBA', startTime: '3:00pm', tv: 'FS1', isSpecial: false, isOffWeek: false, isNonPoints: true },
	{ raceNumber: 13, date: 'Sun., May 24', name: 'Coca-Cola 600', track: 'Charlotte Motor Speedway', stageLaps: '100/200/300/400', startTime: '6:00pm', tv: 'Prime', isSpecial: true, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 14, date: 'Sun., May 31', name: 'Cracker Barrel 400', track: 'Nashville Superspeedway', stageLaps: '90/185/300', startTime: '7:00pm', tv: 'Prime', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 15, date: 'Sun., Jun 7', name: 'FireKeepers Casino 400', track: 'Michigan International Speedway', stageLaps: '45/120/200', startTime: '3:00pm', tv: 'Prime', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 16, date: 'Sun., Jun 14', name: 'Great American Getaway 400', track: 'Pocono Raceway', stageLaps: '30/95/160', startTime: '3:00pm', tv: 'Prime', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 17, date: 'Sun., Jun 21', name: 'Anduril 250 Race the Base', track: 'Naval Base Coronado', stageLaps: 'TBA', startTime: '4:00pm', tv: 'Prime', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 18, date: 'Sun., Jun 28', name: 'Toyota / Save Mart 350', track: 'Sonoma Raceway', stageLaps: '25/55/110', startTime: '3:30pm', tv: 'TNT', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 19, date: 'Sun., Jul 5', name: 'TBA 400', track: 'Chicagoland Speedway', stageLaps: '80/165/267', startTime: '6:00pm', tv: 'TNT', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 20, date: 'Sun., Jul 12', name: 'Quaker State 400', track: 'EchoPark Speedway', stageLaps: '60/160/260', startTime: '7:00pm', tv: 'TNT', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 21, date: 'Sun., Jul 19', name: 'Window World 450', track: 'North Wilkesboro Speedway', stageLaps: 'TBA', startTime: '7:00pm', tv: 'TNT', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 22, date: 'Sun., Jul 26', name: 'Brickyard 400', track: 'Indianapolis Motor Speedway', stageLaps: '50/100/160', startTime: '2:00pm', tv: 'TNT', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: null, date: 'Sun., Aug 2', name: 'Summer Break', track: 'OFF WEEK', stageLaps: '', startTime: '', tv: '', isSpecial: false, isOffWeek: true, isNonPoints: false },
	{ raceNumber: 23, date: 'Sun., Aug 9', name: 'Iowa Corn 350', track: 'Iowa Speedway', stageLaps: '70/210/350', startTime: '3:30pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 24, date: 'Sat., Aug 15', name: 'Cook Out 400', track: 'Richmond Raceway', stageLaps: '70/230/400', startTime: '7:00pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 25, date: 'Sun., Aug 23', name: 'USA Today 301', track: 'New Hampshire Motor Speedway', stageLaps: '70/185/301', startTime: '3:00pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 26, date: 'Sat., Aug 29', name: 'Coke Zero Sugar 400', track: 'Daytona International Speedway', stageLaps: '35/95/160', startTime: '7:30pm', tv: 'NBC', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 27, date: 'Sun., Sep 6', name: 'Cook Out Southern 500', track: 'Darlington Raceway', stageLaps: '115/230/367', startTime: '5:00pm', tv: 'USA', isSpecial: true, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 28, date: 'Sun., Sep 13', name: 'Enjoy Illinois 300', track: 'World Wide Technology Raceway', stageLaps: '45/140/240', startTime: '3:00pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 29, date: 'Sat., Sep 19', name: 'Bass Pro Shops Night Race', track: 'Bristol Motor Speedway', stageLaps: '125/250/500', startTime: '7:30pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 30, date: 'Sun., Sep 27', name: 'Hollywood Casino 400', track: 'Kansas Speedway', stageLaps: '80/165/267', startTime: '3:00pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 31, date: 'Sun., Oct 4', name: 'South Point 400', track: 'Las Vegas Motor Speedway', stageLaps: '80/165/267', startTime: '5:30pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 32, date: 'Sun., Oct 11', name: 'Bank of America ROVAL 400', track: 'Charlotte Motor Speedway ROVAL', stageLaps: '25/50/109', startTime: '3:00pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 33, date: 'Sun., Oct 18', name: 'TBA', track: 'Phoenix Raceway', stageLaps: '60/185/312', startTime: '3:00pm', tv: 'USA', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 34, date: 'Sun., Oct 25', name: 'YellaWood 500', track: 'Talladega Superspeedway', stageLaps: '60/120/188', startTime: '2:00pm', tv: 'NBC', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 35, date: 'Sun., Nov 1', name: 'Xfinity 500', track: 'Martinsville Speedway', stageLaps: '130/260/500', startTime: '2:00pm', tv: 'NBC', isSpecial: false, isOffWeek: false, isNonPoints: false },
	{ raceNumber: 36, date: 'Sun., Nov 8', name: 'NASCAR Cup Series Championship', track: 'Homestead-Miami Speedway', stageLaps: '80/165/267', startTime: '3:00pm', tv: 'NBC', isSpecial: true, isOffWeek: false, isNonPoints: false }
];
