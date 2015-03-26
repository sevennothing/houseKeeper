function genCostIncomeId(keys){
	var id="TT201503251610SJ";
	var now=new Date();
	id = "CIT" + now.getFullYear().toString() + (now.getMonth() + 1).toString()
		+ now.getDate().toString() + now.getHours().toString()
		+ now.getMinutes().toString() + now.getSeconds().toString()
		+ keys;
	return id;
}

function checkVersion(){
	
}
