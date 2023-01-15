package articletable

const (
	template string = `
<el-table ref="articleTable"
        :border=true
        :data="filteredArticle"
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@row-dblclick="HandleDoubleClickedRow"
>
<!--		:default-sort = "{prop: 'Stay.EndDate', order: 'descending'}"-->
	
	<!--	Index   -->
	<el-table-column
			label="N°" width="40px"
			type="index"
			index=1 
	></el-table-column>

	<!--	Actions   -->
	<el-table-column label="Chantiers" width="100px">
		<template slot-scope="scope">
			<div class="header-menu-container">
				<el-button v-if="user.Permissions.Admin" type="danger" plain icon="fa-solid fa-house-circle-xmark fa-fw" size="mini" @click="DeleteAccomodation(scope.row)"></el-button>
				<el-popover v-if="NbSites(scope.row) > 0" placement="right" width="300" trigger="hover" :open-delay=250
					title="Sites associés"
				>
					<div v-for="(binding, index) in BoundSites(scope.row)" :key="binding.Id" style="font-size: 0.85em">
						{{index+1}} - <el-link :href="GetSiteUrl(binding)" target="_blank">{{binding.Site}}</el-link>
					</div>  
					<span slot="reference">{{FormatNbSites(scope.row)}}<i class="fas fa-info-circle icon--right show"></i></span>
				</el-popover>
				<span v-else>-</span>
			</div>
		</template>
	</el-table-column>

	<!--	Ref   -->
	<el-table-column
			:resizable="true" :show-overflow-tooltip=true 
			prop="Ref" label="Nom" width="180px"
	></el-table-column>
	
	<!--	Zip   -->
	<el-table-column
			:resizable="true" :show-overflow-tooltip=true 
			prop="Zipcode" label="Département" width="100px" align="center"
	></el-table-column>
	
	<!--	Address   -->
	<el-table-column
			:resizable="true" :show-overflow-tooltip=true 
			prop="Address" label="Adresse" width="210px"
	>
		<template slot-scope="scope">
			<span>{{FormatAddress(scope.row)}}</span>
		</template>
	</el-table-column>
	
	<!--	Beds & Rooms   -->
	<el-table-column
			:resizable="true" :show-overflow-tooltip=true 
			prop="NumberOfBeds" label="Chambres / Lits" width="120px" align="center"
	>
		<template slot-scope="scope">
			<span>{{FormatBeds(scope.row)}}</span>
		</template>
	</el-table-column>
	
	<!--	Period   -->
	<el-table-column
			:resizable="true" :show-overflow-tooltip=true 
			label="Dates" width="300px"
	>
		<template slot-scope="scope">
			<div  class="header-menu-container on-hover link" @click="EditAccomodationRentalStay(scope.row)">
				<span>{{FormatPeriod(scope.row)}}</span>
				<span>{{FormatDuration(scope.row)}}</span>
			</div>
		</template>
	</el-table-column>

	<!--	Price   -->
	<el-table-column v-if="user.Permissions.Invoice"
			:resizable="true" :show-overflow-tooltip=true 
			label="Prix séjour" width="140px" align="center"
	>
		<template slot-scope="scope">
			<span>{{FormatPrice(scope.row)}}</span>
		</template>
	</el-table-column>

	<!--	Comment   -->
	<el-table-column
			prop="Comment" label="Commentaire"
	></el-table-column>

</el-table>

`
)
