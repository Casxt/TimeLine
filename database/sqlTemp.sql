SELECT `U`.`NickName`, `S`.`Content`, `S`.`Type`, `S`.`Visibility`, `S`.`Location`, `S`.`Time`
					FROM `User` `U`, `Slice` `S`, `Line` `L`
					WHERE `S`.`Visibility`="Public" AND `L`.`Name`=`kkxx` AND `S`.`LineID`=`L`.`ID` AND `U`.`ID`=`S`.`UserID`
					ORDER BY `S`.`Time` DESC LIMIT 0, 20
                    
SELECT `U`.`NickName`, `S`.`Content`, `S`.`Type`, `S`.`Visibility`, `S`.`Location`, `S`.`Time`
					FROM `User` `U`, `Slice` `S`, `Line` `L` 
					WHERE `U`.`ID`=`S`.`UserID` AND `L`.`Name`=`kkxx` AND `S`.`LineID`=`L`.`ID` AND 
					(`S`.`Visibility`=`Public` OR `S`.`Visibility`=`Protect` OR `S`.`UserID`=`116`) 
					ORDER BY `S`.`Time` DESC LIMIT 0, 20;
                    
SELECT `U`.`NickName`, `S`.`Content`, `S`.`Type`, `S`.`Visibility`, `S`.`Location`, `S`.`Time`
					FROM (`User` `U` INNER JOIN `Slice` `S` ON `U`.`ID`=`S`.`UserID`) INNER JOIN `Line` `L` ON `S`.`LineID`=`L`.`ID`
					WHERE `L`.`Name`=`kkxx` AND 
					(`S`.`Visibility`=`Public` OR `S`.`Visibility`=`Protect` OR `S`.`UserID`=`116`) 
					ORDER BY `S`.`Time` DESC LIMIT 0, 20;
SELECT COUNT(*) FROM `Group` `G` INNER JOIN `Line` `L` ON `G`.`LineID`=`L`.`ID` 
					WHERE `L`.`Name`=`kkxx` AND `G`.`UserID`=`116`;
SELECT EXISTS(SELECT * FROM `Group` `G` INNER JOIN `Line` `L` ON `G`.`LineID`=`L`.`ID` 
					WHERE `L`.`Name`='kkxx' AND `G`.`UserID`='116');
SELECT `U`.`NickName`, `S`.`Content`, `S`.`Type`, `S`.`Visibility`, `S`.`Location`, `S`.`Time`
					FROM (`User` `U` INNER JOIN `Slice` `S` ON `U`.`ID`=`S`.`UserID`) INNER JOIN `Line` `L` ON `S`.`LineID`=`L`.`ID`
					WHERE `L`.`Name`='kkxx' AND (`S`.`Visibility` IN ('Public', 'Protect') OR `S`.`UserID`='116') 
					ORDER BY `S`.`Time` DESC LIMIT 0, 20;
                    